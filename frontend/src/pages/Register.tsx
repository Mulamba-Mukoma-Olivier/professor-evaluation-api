import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import {
  authApi,
  RegisterRequest,
} from '../api/auth';

import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';

import {
  GraduationCap,
  ArrowRight,
  ShieldCheck,
  CheckCircle2,
} from 'lucide-react';

export const Register: React.FC = () => {
  const navigate = useNavigate();

  const [formData, setFormData] =
    useState<RegisterRequest>({
      matricule: '',
      name: '',
      email: '',
      password: '',
      role: 'student',
    });

  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (
    e: React.FormEvent
  ) => {
    e.preventDefault();

    setError('');
    setIsLoading(true);

    try {
      await authApi.register(formData);
      navigate('/login');
    } catch (err: any) {
      setError(
        err.response?.data?.error ||
          'Registration failed. Please try again.'
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-[#f4f7fb] p-4 sm:p-6">

      {/* Background decoration */}
      <div className="pointer-events-none absolute inset-0 overflow-hidden">
        <div className="absolute -left-32 -top-32 h-80 w-80 rounded-full bg-blue-500/10 blur-3xl" />
        <div className="absolute -bottom-32 -right-32 h-80 w-80 rounded-full bg-cyan-500/10 blur-3xl" />
      </div>

      <div className="relative grid w-full max-w-5xl overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-2xl shadow-slate-300/40 lg:grid-cols-[1.05fr_0.95fr]">

        {/* =====================================================
            LEFT PANEL
        ===================================================== */}
        <div className="relative hidden min-h-[700px] overflow-hidden bg-[#0b1220] p-8 text-white lg:flex lg:flex-col lg:justify-between">

          {/* Decorative shapes */}
          <div className="pointer-events-none absolute -right-24 -top-24 h-64 w-64 rounded-full bg-blue-600/20 blur-2xl" />

          <div className="pointer-events-none absolute -bottom-24 -left-24 h-64 w-64 rounded-full bg-cyan-500/10 blur-2xl" />

          {/* Logo */}
          <div className="relative flex items-center gap-3">

            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 shadow-lg shadow-blue-900/30">
              <GraduationCap className="h-6 w-6 text-white" />
            </div>

            <div>
              <p className="text-[15px] font-bold tracking-wide">
                UPN-CISNET
              </p>

              <p className="mt-0.5 text-[10px] font-medium uppercase tracking-[0.14em] text-slate-400">
                Evaluation System
              </p>
            </div>

          </div>

          {/* Main message */}
          <div className="relative max-w-md">

            <div className="mb-4 flex items-center gap-2">
              <span className="h-1.5 w-1.5 rounded-full bg-cyan-400" />

              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-400">
                Join CISNET
              </p>
            </div>

            <h1 className="text-4xl font-bold leading-tight tracking-tight">
              Better academic feedback.
              <span className="mt-1 block text-blue-400">
                Better evaluation.
              </span>
            </h1>

            <p className="mt-5 max-w-sm text-sm leading-6 text-slate-400">
              Create your account and access the UPN-CISNET
              platform for structured academic evaluation.
            </p>

            {/* Benefits */}
            <div className="mt-8 space-y-3">

              <div className="flex items-center gap-3 text-sm text-slate-300">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white/5">
                  <CheckCircle2 className="h-4 w-4 text-cyan-400" />
                </div>

                <span>Simple evaluation workflow</span>
              </div>

              <div className="flex items-center gap-3 text-sm text-slate-300">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white/5">
                  <CheckCircle2 className="h-4 w-4 text-cyan-400" />
                </div>

                <span>Centralized academic information</span>
              </div>

              <div className="flex items-center gap-3 text-sm text-slate-300">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-white/5">
                  <CheckCircle2 className="h-4 w-4 text-cyan-400" />
                </div>

                <span>Transparent professor evaluation</span>
              </div>

            </div>
          </div>

          {/* Footer */}
          <div className="relative flex items-center gap-2 text-xs text-slate-500">
            <ShieldCheck className="h-4 w-4 text-blue-400" />

            <span>
              Secure academic evaluation environment
            </span>
          </div>

        </div>

        {/* =====================================================
            RIGHT PANEL
        ===================================================== */}
        <div className="flex min-h-[700px] items-center justify-center bg-white p-6 sm:p-10 lg:p-12">

          <div className="w-full max-w-md">

            {/* Header */}
            <div className="mb-7">

              <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 lg:hidden">
                <GraduationCap className="h-6 w-6" />
              </div>

              <div className="mb-2 flex items-center gap-2">
                <span className="h-1.5 w-1.5 rounded-full bg-blue-600" />

                <span className="text-[11px] font-semibold uppercase tracking-[0.18em] text-blue-600">
                  UPN-CISNET
                </span>
              </div>

              <h2 className="text-3xl font-bold tracking-tight text-slate-900">
                Create account
              </h2>

              <p className="mt-2 text-sm leading-6 text-slate-500">
                Join the UPN-CISNET evaluation system.
              </p>

            </div>

            {/* Form */}
            <form
              onSubmit={handleSubmit}
              className="space-y-4"
            >

              <Input
                label="Matricule"
                type="text"
                value={formData.matricule}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    matricule: e.target.value,
                  })
                }
                placeholder="Enter your matricule"
                required
              />

              <Input
                label="Full Name"
                type="text"
                value={formData.name}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    name: e.target.value,
                  })
                }
                placeholder="Enter your full name"
                required
              />

              <Input
                label="Email"
                type="email"
                value={formData.email}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    email: e.target.value,
                  })
                }
                placeholder="student@upn.edu"
                required
              />

              <Input
                label="Password"
                type="password"
                value={formData.password}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    password: e.target.value,
                  })
                }
                placeholder="Create a password"
                required
              />

              {/* Role */}
              <div className="w-full">

                <label
                  htmlFor="role"
                  className="mb-1.5 block text-sm font-medium text-slate-700"
                >
                  Role
                </label>

                <select
                  id="role"
                  value={formData.role}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      role: e.target.value,
                    })
                  }
                  className="
                    w-full
                    rounded-xl
                    border
                    border-slate-200
                    bg-white
                    px-3.5
                    py-2.5
                    text-sm
                    text-slate-900
                    transition-all
                    duration-200
                    hover:border-slate-300
                    focus:border-blue-500
                    focus:outline-none
                    focus:ring-4
                    focus:ring-blue-100
                  "
                >
                  <option value="student">
                    Student
                  </option>

                  <option value="admin">
                    Admin
                  </option>
                </select>

              </div>

              {/* Error */}
              {error && (
                <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3">
                  <p className="text-sm font-medium leading-5 text-red-700">
                    {error}
                  </p>
                </div>
              )}

              {/* Submit */}
              <Button
                type="submit"
                className="w-full"
                disabled={isLoading}
              >
                <span className="flex items-center justify-center gap-2">

                  {isLoading ? (
                    <>
                      <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />

                      <span>
                        Creating account...
                      </span>
                    </>
                  ) : (
                    <>
                      <span>
                        Create account
                      </span>

                      <ArrowRight className="h-4 w-4" />
                    </>
                  )}

                </span>
              </Button>

            </form>

            {/* Login */}
            <div className="mt-7 border-t border-slate-100 pt-6 text-center">

              <p className="text-sm text-slate-500">
                Already have an account?{' '}

                <button
                  type="button"
                  onClick={() => navigate('/login')}
                  className="font-semibold text-blue-600 transition-colors hover:text-blue-700"
                >
                  Sign in
                </button>
              </p>

            </div>

            {/* Mobile footer */}
            <div className="mt-7 flex items-center justify-center gap-2 text-xs text-slate-400 lg:hidden">
              <ShieldCheck className="h-4 w-4 text-blue-500" />

              <span>
                Secure academic evaluation environment
              </span>
            </div>

          </div>

        </div>

      </div>

    </div>
  );
};