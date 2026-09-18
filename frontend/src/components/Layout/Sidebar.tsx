import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

import {
  LayoutDashboard,
  Users,
  ClipboardList,
  BarChart3,
  LogOut,
  GraduationCap,
  Menu,
  X,
  ChevronRight,
} from 'lucide-react';

interface NavItem {
  label: string;
  icon: React.ReactNode;
  path: string;
}

export const Sidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const navItems: NavItem[] = [
    {
      label: 'Dashboard',
      icon: <LayoutDashboard className="h-[19px] w-[19px]" />,
      path: '/dashboard',
    },
    {
      label: 'Professors',
      icon: <Users className="h-[19px] w-[19px]" />,
      path: '/professors',
    },
    {
      label: 'Evaluations',
      icon: <ClipboardList className="h-[19px] w-[19px]" />,
      path: '/evaluations',
    },
    {
      label: 'Results',
      icon: <BarChart3 className="h-[19px] w-[19px]" />,
      path: '/results',
    },
  ];

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  const handleNavigate = (path: string) => {
    navigate(path);
    setIsMobileMenuOpen(false);
  };

  const handleLogout = () => {
    logout();
    setIsMobileMenuOpen(false);
  };

  return (
    <>
      {/* Mobile menu button */}
      <button
        type="button"
        onClick={() => setIsMobileMenuOpen(true)}
        className="
          fixed left-4 top-4 z-50
          flex h-11 w-11 items-center justify-center
          rounded-xl
          bg-blue-600
          text-white
          shadow-lg shadow-blue-600/20
          transition-all duration-200
          hover:bg-blue-700
          lg:hidden
        "
        aria-label="Open navigation"
      >
        <Menu className="h-5 w-5" />
      </button>

      {/* Mobile overlay */}
      {isMobileMenuOpen && (
        <div
          className="
            fixed inset-0 z-40
            bg-slate-950/50
            backdrop-blur-sm
            lg:hidden
          "
          onClick={() => setIsMobileMenuOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`
          fixed left-0 top-0 z-50
          flex h-screen w-[250px] flex-col
          border-r border-slate-800
          bg-[#0b1220]
          text-white
          shadow-2xl shadow-slate-950/10
          transition-transform duration-300

          ${isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full'}

          lg:static
          lg:translate-x-0
          lg:shadow-none
        `}
      >

        {/* =====================================================
            BRAND
        ====================================================== */}
        <div className="border-b border-slate-800/80 px-5 py-5">
          <div className="flex items-center justify-between">

            <div className="flex items-center gap-3">

              {/* Logo */}
              <div
                className="
                  relative flex h-11 w-11 shrink-0
                  items-center justify-center
                  overflow-hidden rounded-xl
                  bg-gradient-to-br
                  from-blue-500
                  via-blue-600
                  to-cyan-600
                  shadow-lg shadow-blue-600/20
                "
              >
                <GraduationCap className="h-6 w-6 text-white" />

                <div className="absolute -right-3 -top-3 h-7 w-7 rounded-full bg-white/10" />
              </div>

              {/* Brand */}
              <div className="min-w-0">
                <h1 className="text-[15px] font-bold tracking-wide text-white">
                  UPN-CISNET
                </h1>

                <p className="mt-0.5 text-[10px] font-medium uppercase tracking-[0.14em] text-slate-400">
                  Evaluation System
                </p>
              </div>
            </div>

            {/* Close button mobile */}
            <button
              type="button"
              onClick={() => setIsMobileMenuOpen(false)}
              className="
                rounded-lg p-1.5
                text-slate-400
                transition-colors
                hover:bg-slate-800
                hover:text-white
                lg:hidden
              "
              aria-label="Close navigation"
            >
              <X className="h-5 w-5" />
            </button>

          </div>
        </div>

        {/* =====================================================
            NAVIGATION
        ====================================================== */}
        <div className="flex-1 overflow-y-auto px-4 py-6">

          <p className="mb-3 px-3 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
            Main Menu
          </p>

          <nav>
            <ul className="space-y-1.5">

              {navItems.map((item) => {
                const active = isActive(item.path);

                return (
                  <li key={item.path}>
                    <button
                      type="button"
                      onClick={() => handleNavigate(item.path)}
                      className={`
                        group relative flex w-full items-center
                        gap-3 rounded-xl px-3 py-3
                        text-left
                        transition-all duration-200

                        ${
                          active
                            ? `
                              bg-blue-600
                              text-white
                              shadow-lg
                              shadow-blue-600/20
                            `
                            : `
                              text-slate-400
                              hover:bg-slate-800/80
                              hover:text-white
                            `
                        }
                      `}
                    >

                      {/* Active indicator */}
                      {active && (
                        <span
                          className="
                            absolute left-0 top-1/2
                            h-6 w-1
                            -translate-y-1/2
                            rounded-r-full
                            bg-cyan-300
                          "
                        />
                      )}

                      {/* Icon */}
                      <span
                        className={`
                          flex h-9 w-9 shrink-0
                          items-center justify-center
                          rounded-lg
                          transition-colors

                          ${
                            active
                              ? 'bg-white/15 text-white'
                              : 'bg-slate-800/60 text-slate-400 group-hover:text-blue-400'
                          }
                        `}
                      >
                        {item.icon}
                      </span>

                      {/* Label */}
                      <span className="flex-1 text-sm font-medium">
                        {item.label}
                      </span>

                      {/* Arrow */}
                      <ChevronRight
                        className={`
                          h-4 w-4
                          transition-all duration-200

                          ${
                            active
                              ? 'translate-x-0 text-white/80'
                              : '-translate-x-1 text-slate-600 opacity-0 group-hover:translate-x-0 group-hover:opacity-100'
                          }
                        `}
                      />
                    </button>
                  </li>
                );
              })}

            </ul>
          </nav>
        </div>

        {/* =====================================================
            USER / FOOTER
        ====================================================== */}
        <div className="border-t border-slate-800/80 p-4">

          {/* User card */}
          <div
            className="
              mb-3 flex items-center gap-3
              rounded-xl
              border border-slate-800
              bg-slate-900/70
              px-3 py-3
            "
          >

            {/* Avatar */}
            <div
              className="
                flex h-9 w-9 shrink-0
                items-center justify-center
                rounded-lg
                bg-gradient-to-br
                from-blue-500
                to-cyan-600
                text-xs font-bold text-white
              "
            >
              {user?.matricule?.charAt(0)?.toUpperCase() || 'U'}
            </div>

            {/* User information */}
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-semibold text-white">
                {user?.matricule || 'Utilisateur'}
              </p>

              <p className="mt-0.5 truncate text-[10px] font-medium uppercase tracking-wide text-slate-500">
                {user?.role || 'Utilisateur'}
              </p>
            </div>
          </div>

          {/* Logout */}
          <button
            type="button"
            onClick={handleLogout}
            className="
              group flex w-full items-center gap-3
              rounded-xl
              px-3 py-2.5
              text-slate-400
              transition-all duration-200
              hover:bg-red-500/10
              hover:text-red-400
            "
          >
            <span
              className="
                flex h-8 w-8
                items-center justify-center
                rounded-lg
                bg-slate-800
                transition-colors
                group-hover:bg-red-500/10
              "
            >
              <LogOut className="h-4 w-4" />
            </span>

            <span className="text-sm font-medium">
              Logout
            </span>
          </button>

          {/* Footer */}
          <p className="mt-4 text-center text-[9px] uppercase tracking-[0.15em] text-slate-600">
            UPN • CISNET
          </p>

        </div>
      </aside>
    </>
  );
};
