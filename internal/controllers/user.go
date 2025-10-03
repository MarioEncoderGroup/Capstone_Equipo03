package controllers

import (
	userDomain "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	rolePorts "github.com/JoseLuis21/mv-backend/internal/core/role/ports"
	userRoleDomain "github.com/JoseLuis21/mv-backend/internal/core/user_role/domain"
	userRolePorts "github.com/JoseLuis21/mv-backend/internal/core/user_role/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserController maneja las operaciones CRUD de usuarios
type UserController struct {
	userService     ports.UserService
	roleService     rolePorts.RoleService
	userRoleService userRolePorts.UserRoleService
	validator       *validatorapi.Validator
}

// NewUserController crea una nueva instancia del controller de usuarios
func NewUserController(userService ports.UserService, roleService rolePorts.RoleService, userRoleService userRolePorts.UserRoleService, validator *validatorapi.Validator) *UserController {
	return &UserController{
		userService:     userService,
		roleService:     roleService,
		userRoleService: userRoleService,
		validator:       validator,
	}
}

// GetUsers maneja GET /users - Lista usuarios con paginación
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	// Parsear parámetros de paginación
	paginationReq, err := pagination.ParsePaginationFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Parámetros de paginación inválidos",
			Error:   err.Error(),
		})
	}

	// Validar campo de ordenamiento para usuarios
	allowedSortFields := []string{"id", "username", "full_name", "email", "created", "updated", "is_active"}
	if err := paginationReq.SetCustomSortField(paginationReq.SortBy, allowedSortFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Campo de ordenamiento inválido",
			Error:   err.Error(),
		})
	}

	// Obtener usuarios desde el servicio
	users, total, err := uc.userService.GetUsers(
		c.Context(),
		paginationReq.GetOffset(),
		paginationReq.GetLimit(),
		paginationReq.SortBy,
		paginationReq.SortDir,
		paginationReq.Search,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo usuarios",
			Error:   err.Error(),
		})
	}

	// Convertir a DTOs de respuesta y cargar roles
	var userDtos []userDomain.UserListResponseDto
	for _, user := range users {
		dto := user.ToUserListResponseDto()
		
		// Cargar roles del usuario (sin filtro de tenant para obtener todos)
		roles, err := uc.roleService.GetUserRoles(c.Context(), user.ID, nil)
		if err == nil && len(roles) > 0 {
			dto.Roles = make([]userDomain.RoleResponseDto, len(roles))
			for i, role := range roles {
				dto.Roles[i] = userDomain.RoleResponseDto{
					ID:          role.ID,
					Name:        role.Name,
					Description: role.Description,
				}
			}
		}
		
		userDtos = append(userDtos, *dto)
	}

	// Calcular información de paginación
	paginationInfo := paginationReq.CalculatePagination(total)

	return c.Status(fiber.StatusOK).JSON(types.PaginatedAPIResponse{
		Success:    true,
		Message:    "Usuarios obtenidos exitosamente",
		Data:       userDtos,
		Pagination: paginationInfo,
	})
}

// GetUserByID maneja GET /users/:id - Obtiene un usuario específico
func (uc *UserController) GetUserByID(c *fiber.Ctx) error {
	// Parsear ID del usuario
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Obtener usuario desde el servicio
	user, err := uc.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no encontrado",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario obtenido exitosamente",
		Data:    user.ToUserResponseDto(),
	})
}

// CreateUser maneja POST /users - Crea un nuevo usuario
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	// Parsear request body
	var req userDomain.CreateUserDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := uc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Crear usuario a través del servicio
	user, err := uc.userService.CreateUserFromDto(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error creando usuario",
			Error:   err.Error(),
		})
	}

	// Asignar roles si se enviaron en el request
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			createUserRoleDto := &userRoleDomain.CreateUserRoleDto{
				UserID:   user.ID,
				RoleID:   roleID,
				TenantID: nil, // Sin tenant específico (rol global)
			}

			_, err := uc.userRoleService.CreateUserRole(c.Context(), createUserRoleDto)
			if err != nil {
				// Log error pero no fallar la creación del usuario
				// El usuario ya fue creado, solo falla la asignación del rol
				// Podrías decidir retornar error o solo loguearlo
			}
		}
	}

	// Preparar respuesta con los roles asignados
	dto := user.ToUserListResponseDto()
	
	// Cargar roles del usuario recién creado
	roles, err := uc.roleService.GetUserRoles(c.Context(), user.ID, nil)
	if err == nil && len(roles) > 0 {
		dto.Roles = make([]userDomain.RoleResponseDto, len(roles))
		for i, role := range roles {
			dto.Roles[i] = userDomain.RoleResponseDto{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario creado exitosamente",
		Data:    dto,
	})
}

// UpdateUser maneja PUT /users/:id - Actualiza un usuario
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	// Parsear ID del usuario
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var req userDomain.UpdateUserDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := uc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Actualizar usuario a través del servicio
	user, err := uc.userService.UpdateUserFromDto(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario actualizado exitosamente",
		Data:    user.ToUserResponseDto(),
	})
}

// DeleteUser maneja DELETE /users/:id - Elimina (soft delete) un usuario
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	// Parsear ID del usuario
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Eliminar usuario a través del servicio
	if err := uc.userService.DeleteUser(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error eliminando usuario",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario eliminado exitosamente",
	})
}

// ChangePassword maneja POST /users/:id/change-password - Cambia la contraseña de un usuario
func (uc *UserController) ChangePassword(c *fiber.Ctx) error {
	// Parsear ID del usuario
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido",
			Error:   "El ID debe ser un UUID válido",
		})
	}

	// Parsear request body
	var req userDomain.ChangePasswordDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := uc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Cambiar contraseña a través del servicio
	if err := uc.userService.ChangeUserPassword(c.Context(), userID, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error cambiando contraseña",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Contraseña actualizada exitosamente",
	})
}

// GetProfile maneja GET /users/profile - Obtiene el perfil del usuario autenticado
func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	// Extraer ID del usuario autenticado del contexto
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "Token de autenticación requerido",
		})
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido en token",
			Error:   "Token corrupto",
		})
	}

	// Obtener usuario desde el servicio
	user, err := uc.userService.GetUserByID(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no encontrado",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Perfil obtenido exitosamente",
		Data:    user.ToUserResponseDto(),
	})
}

// UpdateProfile maneja PUT /users/profile - Actualiza el perfil del usuario autenticado
func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	// Extraer ID del usuario autenticado del contexto
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.APIResponse{
			Success: false,
			Message: "Usuario no autenticado",
			Error:   "Token de autenticación requerido",
		})
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de usuario inválido en token",
			Error:   "Token corrupto",
		})
	}

	// Parsear request body
	var req userDomain.UpdateProfileDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error parseando datos de entrada",
			Error:   "Formato de datos inválido",
		})
	}

	// Validar estructura de datos
	if errors := uc.validator.ValidateStruct(req); len(errors) > 0 {
		var validationErrors []types.ValidationErrorResponse
		for _, err := range errors {
			validationErrors = append(validationErrors, types.ValidationErrorResponse{
				Field:   err.Field,
				Message: err.Message,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Datos de entrada inválidos",
			Data:    validationErrors,
		})
	}

	// Actualizar perfil a través del servicio
	user, err := uc.userService.UpdateUserProfile(c.Context(), userUUID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "Error actualizando perfil",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Perfil actualizado exitosamente",
		Data:    user.ToUserResponseDto(),
	})
}
