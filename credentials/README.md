# Credenciales

Este directorio contiene las credenciales de Google Cloud Vision API.

## Setup

1. Descargar el archivo JSON de credenciales desde Google Cloud Console
2. Guardarlo aquí con el nombre: `google-vision-key.json`
3. Configurar la variable de entorno en `.env`:

```bash
GOOGLE_APPLICATION_CREDENTIALS=./credentials/google-vision-key.json
```

## Seguridad

⚠️ **IMPORTANTE**: 
- Este directorio está en `.gitignore`
- NUNCA commitear archivos JSON de credenciales
- Rotar credenciales cada 90 días
- No compartir credenciales por email/slack

## Verificación

Para verificar que las credenciales están configuradas:

```bash
# Verificar que el archivo existe
ls -la credentials/google-vision-key.json

# Verificar que NO está siendo tracked por git
git status  # No debe aparecer en la lista
```
