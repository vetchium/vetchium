# Vetchi Jobs - Hub UI

This is the Hub UI application for Vetchi Jobs, where job seekers can search for openings, submit applications, and track their candidacies.

## Tech Stack

- React 18
- TypeScript
- Vite
- Ant Design (UI Components)
- Redux Toolkit (State Management)
- Axios (API Client)
- Styled Components (CSS-in-JS)

## Prerequisites

- Node.js (v16 or later)
- npm or yarn

## Getting Started

1. Install dependencies:
```bash
npm install
# or
yarn
```

2. Set up environment variables:
```bash
cp .env.example .env
```
Edit `.env` file and update the API endpoint if needed.

3. Start development server:
```bash
npm run dev
# or
yarn dev
```

The application will be available at `http://localhost:3001`.

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

## Project Structure

```
hub-ui/
├── src/
│   ├── components/     # Reusable UI components
│   ├── hooks/         # Custom React hooks
│   ├── layouts/       # Layout components
│   ├── pages/         # Page components
│   ├── store/         # Redux store and slices
│   ├── types/         # TypeScript type definitions
│   ├── App.tsx        # Root component
│   └── main.tsx       # Entry point
├── public/            # Static assets
└── index.html         # HTML template
```

## Key Features

1. Authentication
   - Sign in with email and password
   - Automatic token refresh
   - Protected routes

2. Job Search
   - Filter by employer, location, job type
   - Country-based location filtering
   - Application submission with resume and cover letter

3. Application Management
   - Track application status
   - Communicate with employers
   - Withdraw applications

4. Candidacy Tracking
   - View interview schedule
   - Track candidacy progress
   - Respond to job offers

5. Profile Management
   - Update personal information
   - Upload/update resume
   - Change password
   - Theme and language preferences

## Development Guidelines

1. Code Style
   - Follow ESLint and Prettier configurations
   - Use TypeScript strict mode
   - Write meaningful component and variable names

2. Component Structure
   - Keep components focused and single-responsibility
   - Use TypeScript interfaces for props
   - Implement proper error handling
   - Add loading states for async operations

3. State Management
   - Use Redux for global state
   - Use local state for component-specific data
   - Implement proper error handling in API calls

4. Testing
   - Write unit tests for components
   - Test error scenarios
   - Verify form validations

5. Making Changes
   - Create a new branch for features/fixes
   - Follow commit message conventions
   - Update documentation when needed
   - Test thoroughly before submitting PR

## API Integration

The application communicates with the backend API at `example.com`. All API endpoints are prefixed with `/api/hub/`.

Common API patterns:
- Use axios interceptors for token handling
- Implement proper error handling
- Show loading states during API calls
- Display success/error messages to users

## Deployment

1. Build the application:
```bash
npm run build
# or
yarn build
```

2. The built files will be in the `dist` directory, ready to be deployed.

3. For production deployment:
   - Set up proper environment variables
   - Configure proper API endpoints
   - Set up proper CORS settings
   - Configure proper security headers

## Troubleshooting

1. Common Issues
   - CORS errors: Check API endpoint configuration
   - Build errors: Check dependencies and TypeScript errors
   - Runtime errors: Check console logs and error boundaries

2. Development Tips
   - Use React DevTools for component debugging
   - Use Redux DevTools for state debugging
   - Check browser console for errors
   - Verify API responses in Network tab

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is proprietary and confidential. 