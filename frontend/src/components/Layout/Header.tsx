import React from 'react';
import { Bell, ChevronDown, User } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

export const Header: React.FC<{ title: string }> = ({ title }) => {
  const { user } = useAuth();

  return (
    <header className="sticky top-0 z-30 border-b border-slate-200 bg-white/95 backdrop-blur-xl">
      <div className="flex h-20 items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">

        {/* Left side */}
        <div className="min-w-0">
          <div className="mb-1 flex items-center gap-2">
            <span className="h-1.5 w-1.5 rounded-full bg-blue-600" />

            <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-blue-600">
              CISNET
            </p>
          </div>

          <h1 className="truncate text-xl font-bold tracking-tight text-slate-900 lg:text-2xl">
            {title}
          </h1>
        </div>

        {/* Right side */}
        <div className="flex items-center gap-3">

          {/* Notification */}
          <button
            type="button"
            className="
              relative flex h-10 w-10 items-center justify-center
              rounded-xl border border-slate-200
              bg-white text-slate-500
              transition-all duration-200
              hover:border-blue-200
              hover:bg-blue-50
              hover:text-blue-600
            "
            aria-label="Notifications"
          >
            <Bell className="h-[18px] w-[18px]" />

            <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-blue-600 ring-2 ring-white" />
          </button>

          {/* User profile */}
          <div
            className="
              flex items-center gap-3
              rounded-xl border border-slate-200
              bg-slate-50/70
              px-2.5 py-2
              transition-colors
              hover:border-blue-200
              hover:bg-blue-50/40
            "
          >
            {/* Avatar */}
            <div
              className="
                flex h-10 w-10 shrink-0 items-center justify-center
                rounded-xl
                bg-gradient-to-br from-blue-600 to-cyan-600
                text-white
                shadow-sm
              "
            >
              <User className="h-[18px] w-[18px]" />
            </div>

            {/* User information */}
            <div className="hidden min-w-0 sm:block">
              <p className="max-w-[150px] truncate text-sm font-semibold text-slate-900">
                {user?.matricule || 'Utilisateur'}
              </p>

              <p className="text-[11px] font-medium uppercase tracking-wide text-slate-500">
                {user?.role || 'Utilisateur'}
              </p>
            </div>

            {/* Dropdown indicator */}
            <ChevronDown className="hidden h-4 w-4 text-slate-400 sm:block" />
          </div>
        </div>
      </div>
    </header>
  );
};
