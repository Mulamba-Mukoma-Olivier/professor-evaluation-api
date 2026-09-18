import React from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Header } from './Header';

export const DashboardLayout: React.FC = () => {
  const location = useLocation();

  const getTitle = () => {
    const path = location.pathname;

    if (path === '/dashboard') return 'Dashboard';
    if (path === '/professors') return 'Professors';
    if (path === '/evaluations') return 'Evaluations';
    if (path === '/results') return 'Results';

    return 'Dashboard';
  };

  return (
    <div className="min-h-screen bg-[#f4f7fb] text-slate-800">
      <div className="flex min-h-screen">

        {/* Sidebar */}
        <aside className="hidden lg:block w-64 shrink-0">
          <Sidebar />
        </aside>

        {/* Main application */}
        <div className="flex min-w-0 flex-1 flex-col">

          {/* Header */}
          <header className="sticky top-0 z-40 border-b border-slate-200/80 bg-white/95 backdrop-blur">
            <Header title={getTitle()} />
          </header>

          {/* Content */}
          <main className="flex-1 p-4 sm:p-5 lg:p-7">
            <div className="mx-auto w-full max-w-[1600px]">
              <Outlet />
            </div>
          </main>

          {/* Footer */}
          <footer className="border-t border-slate-200 bg-white px-6 py-4">
            <div className="flex flex-col items-center justify-between gap-2 text-xs text-slate-500 sm:flex-row">
              <span>
                © {new Date().getFullYear()} CISNET
              </span>

              <span>
                Academic Evaluation Management System
              </span>
            </div>
          </footer>

        </div>
      </div>
    </div>
  );
};