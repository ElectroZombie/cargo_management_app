# Cargo Management App

Una aplicación de escritorio moderna para gestión de carga utilizando **Wails**, **Go** y **Angular**.

## 📋 Características

- ✨ Interfaz moderna y responsiva con Angular
- 🚀 Backend robusto con Go
- 📦 Gestión completa de carga
- 📊 Dashboard con estadísticas
- 🎨 Diseño profesional con SCSS

## 🛠️ Requisitos

Antes de comenzar, asegúrate de tener instalado:

- **Go** 1.18+
- **Node.js** 16+ y npm
- **Wails CLI** v2+

Instala Wails globalmente:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## 📁 Estructura del Proyecto

```
cargo_management_app/
├── main.go                 # Punto de entrada de la aplicación
├── wails.json             # Configuración de Wails
├── go.mod                 # Módulo de Go
├── go.sum                 # Dependencias de Go
└── frontend/              # Aplicación Angular
    ├── src/
    │   ├── app/           # Componentes y servicios
    │   ├── index.html     # HTML principal
    │   ├── main.ts        # Bootstrap de Angular
    │   └── styles.scss    # Estilos globales
    ├── package.json       # Dependencias npm
    ├── angular.json       # Configuración de Angular
    └── tsconfig.json      # Configuración de TypeScript
```

## 🚀 Inicio Rápido

### 1. Instalar dependencias

**Backend:**
```bash
go mod download
```

**Frontend:**
```bash
cd frontend
npm install
```

### 2. Desarrollo

Ejecuta la aplicación en modo desarrollo:

```bash
wails dev
```

Esto iniciará:
- La aplicación Wails con hot-reload
- El servidor de desarrollo de Angular
- El backend de Go

### 3. Build para producción

Para compilar la aplicación:

```bash
wails build
```

El ejecutable se encontrará en el directorio `/build`.

## 🏗️ Estructura del Código

### Backend (Go)

El archivo `main.go` contiene:
- Configuración de Wails
- Inicialización de la aplicación
- Métodos expuestos al frontend

### Frontend (Angular)

#### Componentes principales:

- **AppComponent**: Componente raíz
- **NavbarComponent**: Barra de navegación superior
- **SidebarComponent**: Menú lateral
- **DashboardComponent**: Página principal con estadísticas
- **CargoListComponent**: Gestión de carga

#### Rutas:

- `/dashboard` - Panel principal
- `/cargo` - Gestión de carga

## 🎨 Estilos

Se utiliza **SCSS** para los estilos. Los archivos de estilo se encuentran en:

- `frontend/src/styles.scss` - Estilos globales
- `frontend/src/app/**/*.scss` - Estilos de componentes

### Variables y Temas

**Colores principales:**
- Fondo: `#1b2636`
- Navbar: `#2c3e50`
- Sidebar: `#34495e`
- Primario: `#4CAF50`
- Secundario: `#008CBA`

## 📦 Gestión de Cargo

La aplicación permite:

- ✅ Ver lista de carga
- ✅ Agregar nueva carga
- ✅ Actualizar estado
- ✅ Eliminar carga
- ✅ Filtrar por estado

## 🔗 API del Backend

El backend expone métodos Go que pueden ser llamados desde Angular:

```go
// Ejemplo en main.go
func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s, it's been fun!", name)
}
```

Para llamar desde Angular:

```typescript
// Necesitas implementar el servicio de Wails
import { invoke } from '@wailsapp/runtime';

invoke('App.Greet', { name: 'World' })
    .then(result => console.log(result))
    .catch(err => console.error(err));
```

## 🐛 Troubleshooting

### El servidor de desarrollo no inicia

1. Elimina las carpetas: `node_modules/`, `dist/`
2. Reinstala dependencias: `npm install`
3. Intenta nuevamente: `wails dev`

### Errores de módulos Go

```bash
go mod tidy
go mod download
```

### Puerto ya en uso

Por defecto, Wails usa el puerto 34115. Si está en uso:

```bash
wails dev --port 3000
```

## 📚 Recursos Útiles

- [Documentación de Wails](https://wails.io/docs)
- [Documentación de Angular](https://angular.io/docs)
- [Guía de Go](https://golang.org/doc/)

## 📝 Licencia

MIT - Siéntete libre de usar este proyecto como desees.

## 👤 Autor

**ElectroZombie**

---

Creado con ❤️ para gestión eficiente de carga
