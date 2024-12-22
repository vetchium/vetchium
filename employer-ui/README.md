# Vetchi Jobs - Employer UI

This is the Employer UI application for Vetchi Jobs, where employers can manage job openings, review applications, and conduct hiring processes.

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

The application will be available at `http://localhost:3000`.

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run type-check` - Run TypeScript type checking

## Project Structure

```
employer-ui/
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
   - Sign in with domain, email, and password
   - Role-based access control
   - Protected routes

2. Organization Management
   - Manage OrgUsers (add/remove users)
   - Configure Locations
   - Configure Departments
   - Role assignments

3. Job Opening Management
   - Create and edit job openings
   - Assign hiring managers and recruiters
   - Set requirements and qualifications
   - Manage opening status

4. Application Review
   - Review submitted applications
   - Shortlist or reject candidates
   - View resumes and cover letters
   - Communicate with candidates

5. Interview Process
   - Schedule interviews
   - Assign interviewers
   - Collect interview feedback
   - Make hiring decisions

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

The application communicates with the backend API at `example.com`. All API endpoints are prefixed with `/api/employer/`.

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

## Role-Based Access

The application supports different user roles:
- Admin: Full access to all features
- Hiring Manager: Manage openings and make hiring decisions
- Recruiter: Review applications and coordinate interviews
- Interviewer: Provide interview feedback

Ensure proper role checks are implemented for protected features.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is available under GNU Affero General Public License v3.0.