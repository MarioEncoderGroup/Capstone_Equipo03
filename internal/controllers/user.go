package controllers

import (
	userDomain "github.com/JoseLuis21/mv-backend/internal/core/user/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/user/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/pagination"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserController maneja las operaciones CRUD de usuarios
type UserController struct {
	userService ports.UserService
	validator   *validatorapi.Validator
}

// NewUserController crea una nueva instancia del controller de usuarios
func NewUserController(userService ports.UserService, validator *validatorapi.Validator) *UserController {
	return &UserController{
		userService: userService,
		validator:   validator,
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

	// Convertir a DTOs de respuesta
	var userDtos []userDomain.UserListResponseDto
	for _, user := range users {
		userDtos = append(userDtos, *user.ToUserListResponseDto())
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

	return c.Status(fiber.StatusCreated).JSON(types.APIResponse{
		Success: true,
		Message: "Usuario creado exitosamente",
		Data:    user.ToUserResponseDto(),
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
