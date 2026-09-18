# UPN-CISNET Professor Evaluation Frontend

A modern React + TypeScript frontend for the Professor Evaluation System with an Apple-inspired design and enterprise styling.

## Tech Stack

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **TailwindCSS** - Styling with Apple-inspired theme
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **Lucide React** - Icon library

## Features

- **Authentication**: Login and registration with JWT tokens
- **Dashboard**: Overview with statistics and recent activity
- **Professors Management**: View and search professors
- **Evaluations**: Submit and view professor evaluations
- **Results & Analytics**: Detailed professor performance analytics
- **Responsive Design**: Mobile-friendly interface
- **Apple-inspired UI**: Clean, modern design with enterprise styling

## Prerequisites

- Node.js 18+ (Note: Vite 8+ requires Node 20+, but this project is configured for Node 18)
- npm or yarn

## Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

## Configuration

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080
```

## Development

Start the development server:

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

## Build for Production

Build the application:

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── api/              # API service layer
│   │   ├── auth.ts
│   │   ├── professors.ts
│   │   ├── evaluations.ts
│   │   ├── courses.ts
│   │   ├── criteria.ts
│   │   ├── results.ts
│   │   └── client.ts
│   ├── components/
│   │   ├── ui/          # Reusable UI components
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   └── Card.tsx
│   │   └── Layout/      # Layout components
│   │       ├── Sidebar.tsx
│   │       ├── Header.tsx
│   │       └── DashboardLayout.tsx
│   ├── context/         # React context
│   │   └── AuthContext.tsx
│   ├── pages/           # Page components
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Professors.tsx
│   │   ├── Evaluations.tsx
│   │   └── Results.tsx
│   ├── App.tsx          # Main app component with routing
│   ├── main.tsx         # Entry point
│   └── index.css        # Global styles
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
└── package.json
```

## Design System

### Apple-Inspired Colors

- **Primary Blue**: `#007AFF`
- **Gray Scale**: Various shades from light to dark
- **Success Green**: `#34C759`
- **Warning Orange**: `#FF9500`
- **Error Red**: `#FF3B30`

### Enterprise Colors

- **Primary**: `#1a365d`
- **Secondary**: `#2c5282`
- **Accent**: `#3182ce`

### Components

- **Buttons**: Multiple variants (primary, secondary, outline, ghost)
- **Inputs**: Clean, focused design with error states
- **Cards**: Shadowed containers with headers and content
- **Layout**: Fixed sidebar with responsive header

## API Integration

The frontend connects to the Go backend API at the configured `VITE_API_URL`. All API calls include JWT authentication tokens automatically.

### Endpoints Used

- `POST /auth/login` - User authentication
- `POST /auth/register` - User registration
- `GET /professors` - List all professors
- `GET /professors/active` - List active professors
- `GET /evaluations` - List evaluations
- `POST /evaluations/:student_id` - Create evaluation
- `GET /results/professors/:professor_id` - Get professor results
- `GET /courses` - List courses
- `GET /criteria/active` - List active criteria

## Authentication

The application uses JWT tokens for authentication. Tokens are stored in localStorage and automatically included in API requests. Protected routes redirect to login if not authenticated.

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

ISC
