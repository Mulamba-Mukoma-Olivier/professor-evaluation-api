import React from 'react';
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from 'react-router-dom';

import {
  AuthProvider,
  useAuth,
} from './context/AuthContext';

import { Login } from './pages/Login';
import { Register } from './pages/Register';

import {
  DashboardLayout,
} from './components/Layout/DashboardLayout';

import { Dashboard } from './pages/Dashboard';
import { Professors } from './pages/Professors';
import { Evaluations } from './pages/Evaluations';
import { Results } from './pages/Results';

/* =========================================================
   PROTECTED ROUTE
========================================================= */

const ProtectedRoute: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const {
    isAuthenticated,
    isLoading,
  } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#f4f7fb]">
        <div className="flex flex-col items-center gap-3">

          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-blue-100 border-t-blue-600" />
          </div>

          <p className="text-sm font-medium text-slate-500">
            Loading CISNET...
          </p>

        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

/* =========================================================
   PUBLIC ROUTE
========================================================= */

const PublicRoute: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const {
    isAuthenticated,
    isLoading,
  } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#f4f7fb]">
        <div className="flex flex-col items-center gap-3">

          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-blue-100 border-t-blue-600" />
          </div>

          <p className="text-sm font-medium text-slate-500">
            Loading CISNET...
          </p>

        </div>
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
};

/* =========================================================
   APPLICATION ROUTES
========================================================= */

const AppRoutes: React.FC = () => {
  return (
    <Routes>

      {/* ===================================================
          PUBLIC ROUTES
      =================================================== */}

      <Route
        path="/login"
        element={
          <PublicRoute>
            <Login />
          </PublicRoute>
        }
      />

      <Route
        path="/register"
        element={
          <PublicRoute>
            <Register />
          </PublicRoute>
        }
      />

      {/* ===================================================
          PROTECTED APPLICATION
      =================================================== */}

      <Route
        path="/"
        element={
          <ProtectedRoute>
            <DashboardLayout />
          </ProtectedRoute>
        }
      >

        {/* / */}
        <Route
          index
          element={
            <Navigate
              to="/dashboard"
              replace
            />
          }
        />

        {/* /dashboard */}
        <Route
          path="dashboard"
          element={<Dashboard />}
        />

        {/* /professors */}
        <Route
          path="professors"
          element={<Professors />}
        />

        {/* /evaluations */}
        <Route
          path="evaluations"
          element={<Evaluations />}
        />

        {/* /results */}
        <Route
          path="results"
          element={<Results />}
        />

      </Route>

      {/* ===================================================
          FALLBACK
      =================================================== */}

      <Route
        path="*"
        element={
          <Navigate
            to="/dashboard"
            replace
          />
        }
      />

    </Routes>
  );
};

/* =========================================================
   APPLICATION
========================================================= */

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;