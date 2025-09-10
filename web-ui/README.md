# Crypto Inventory Frontend

A modern React-based frontend for the Crypto Inventory Management System, built with Vite, TypeScript, and TailwindCSS.

## 🎨 Design System

### Color Scheme
The application uses a professional Black, Gold, and Red color palette:

- **Primary (Gold)**: Rich gold tones for main actions and branding
  - Main: `#f59e0b` (primary-500)
  - Hover: `#d97706` (primary-600)
  - Light: `#fef3c7` (primary-100)

- **Secondary (Black)**: Deep blacks and grays for backgrounds and text
  - Deep Black: `#0f172a` (secondary-900)
  - Medium: `#334155` (secondary-700)
  - Light: `#f1f5f9` (secondary-100)

- **Accent (Red)**: Vibrant red for highlights, alerts, and important actions
  - Main: `#ef4444` (accent-500)
  - Hover: `#dc2626` (accent-600)
  - Light: `#fef2f2` (accent-50)

### Dark Mode Support
The application includes full dark mode support with automatic theme switching based on user preference.

## 🚀 Development

### Prerequisites
- Node.js 18+
- npm or yarn

### Setup
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Development Server
The development server runs on `http://localhost:5173` by default. If port 5173 is busy, Vite will automatically try the next available port (5174, 5175, etc.).

### Available Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run test` - Run tests
- `npm run test:watch` - Run tests in watch mode

## 🏗️ Architecture

### Tech Stack
- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **TailwindCSS** - Utility-first CSS framework
- **React Router** - Client-side routing
- **TanStack Query** - Data fetching and caching
- **Axios** - HTTP client
- **React Hook Form** - Form handling
- **Zod** - Schema validation
- **Headless UI** - Accessible UI components
- **Heroicons** - Icon library

### Project Structure
```
src/
├── components/     # Reusable UI components
├── pages/         # Page components
├── services/      # API services
├── hooks/         # Custom React hooks
├── utils/         # Utility functions
├── types/         # TypeScript type definitions
└── styles/        # Global styles and CSS
```

## 🎨 Customization

### Changing Colors
The color scheme can be easily customized by updating `tailwind.config.js`:

```javascript
colors: {
  primary: {
    // Update these values to change the primary color
    500: '#f59e0b',  // Main gold
    // ... other shades
  },
  secondary: {
    // Update these values to change the secondary color
    900: '#0f172a',  // Deep black
    // ... other shades
  },
  accent: {
    // Update these values to change the accent color
    500: '#ef4444',  // Main red
    // ... other shades
  },
}
```

### Adding New Components
Use the predefined component classes for consistent styling:

```jsx
// Primary button
<button className="btn-primary">Click me</button>

// Secondary button
<button className="btn-secondary">Cancel</button>

// Accent button
<button className="btn-accent">Delete</button>

// Form input
<input className="form-input" placeholder="Enter text..." />

// Card container
<div className="card">
  <h3>Card Title</h3>
  <p>Card content</p>
</div>
```

## 🔧 Configuration

### Environment Variables
Create a `.env.local` file for local development:

```env
VITE_API_BASE_URL=http://localhost:8081/api/v1
VITE_APP_NAME=Crypto Inventory
VITE_APP_VERSION=1.0.0
```

### TailwindCSS Configuration
The TailwindCSS configuration is located in `tailwind.config.js` and includes:
- Custom color palette
- Dark mode support
- Custom font family
- Component classes

## 📱 Responsive Design
The application is fully responsive and works on:
- Desktop (1024px+)
- Tablet (768px - 1023px)
- Mobile (320px - 767px)

## 🧪 Testing
The project includes comprehensive testing setup with:
- Vitest for unit testing
- React Testing Library for component testing
- Jest DOM for DOM assertions

Run tests with:
```bash
npm run test
npm run test:watch
```

## 🚀 Deployment
The application builds to static files that can be served by any web server:

```bash
npm run build
```

The built files will be in the `dist/` directory and can be served by nginx, Apache, or any static file server.

## 📚 Additional Resources
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
