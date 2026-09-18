import React, { useEffect, useState } from 'react';

import {
  professorsApi,
  Professor,
} from '../api/professors';

import {
  Card,
  CardContent,
} from '../components/ui/Card';

import {
  Button,
} from '../components/ui/Button';

import {
  Input,
} from '../components/ui/Input';

import {
  Search,
  Plus,
  Mail,
  Building2,
  User,
  Users,
  UserCheck,
  UserX,
  ArrowUpRight,
} from 'lucide-react';

export const Professors: React.FC = () => {
  const [professors, setProfessors] = useState<Professor[]>([]);
  const [filteredProfessors, setFilteredProfessors] = useState<
    Professor[]
  >([]);

  const [searchTerm, setSearchTerm] = useState('');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadProfessors();
  }, []);

  useEffect(() => {
    const term = searchTerm.toLowerCase().trim();

    const filtered = professors.filter((professor) => {
      return (
        professor.first_name
          ?.toLowerCase()
          .includes(term) ||
        professor.last_name
          ?.toLowerCase()
          .includes(term) ||
        professor.email
          ?.toLowerCase()
          .includes(term) ||
        professor.department
          ?.toLowerCase()
          .includes(term) ||
        professor.grade
          ?.toLowerCase()
          .includes(term)
      );
    });

    setFilteredProfessors(filtered);
  }, [searchTerm, professors]);

  const loadProfessors = async () => {
    try {
      const data = await professorsApi.getAll();

      setProfessors(data);
      setFilteredProfessors(data);
    } catch (error) {
      console.error(
        'Failed to load professors:',
        error
      );
    } finally {
      setIsLoading(false);
    }
  };

  const activeProfessors = professors.filter(
    (professor) => professor.active
  ).length;

  const inactiveProfessors =
    professors.length - activeProfessors;

  const getInitials = (professor: Professor) => {
    const first = professor.first_name?.charAt(0) || '';
    const last = professor.last_name?.charAt(0) || '';

    return (
      `${first}${last}`.toUpperCase() || '?'
    );
  };

  return (
    <div className="space-y-7">

      {/* =====================================================
          HEADER
      ===================================================== */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">

        <div>
          <div className="mb-1 flex items-center gap-2">
            <span className="h-1.5 w-1.5 rounded-full bg-blue-600" />

            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-blue-600">
              Academic Staff
            </p>
          </div>

          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Professors
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Manage and view professors registered in CISNET.
          </p>
        </div>

        <Button
          type="button"
          className="w-full sm:w-auto"
        >
          <Plus className="h-4 w-4" />
          <span>Add Professor</span>
        </Button>

      </div>

      {/* =====================================================
          SUMMARY CARDS
      ===================================================== */}
      {!isLoading && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">

          {/* Total */}
          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">

                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Total Professors
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {professors.length}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                  <Users className="h-5 w-5" />
                </div>

              </div>
            </CardContent>
          </Card>

          {/* Active */}
          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">

                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Active Professors
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {activeProfessors}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                  <UserCheck className="h-5 w-5" />
                </div>

              </div>
            </CardContent>
          </Card>

          {/* Inactive */}
          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">

                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Inactive Professors
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {inactiveProfessors}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-slate-100 text-slate-500">
                  <UserX className="h-5 w-5" />
                </div>

              </div>
            </CardContent>
          </Card>

        </div>
      )}

      {/* =====================================================
          SEARCH
      ===================================================== */}
      <Card>
        <CardContent className="p-4 sm:p-5">

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">

            <div className="relative w-full sm:max-w-md">

              <Search
                className="
                  absolute
                  left-3
                  top-1/2
                  h-4
                  w-4
                  -translate-y-1/2
                  text-slate-400
                "
              />

              <Input
                placeholder="Search by name, email, department..."
                value={searchTerm}
                onChange={(e) =>
                  setSearchTerm(e.target.value)
                }
                className="pl-10"
              />

            </div>

            <div className="text-xs font-medium text-slate-500">
              {searchTerm
                ? `${filteredProfessors.length} result${
                    filteredProfessors.length !== 1
                      ? 's'
                      : ''
                  }`
                : `${professors.length} professors`}
            </div>

          </div>

        </CardContent>
      </Card>

      {/* =====================================================
          LOADING
      ===================================================== */}
      {isLoading ? (
        <Card>
          <CardContent className="flex min-h-[300px] items-center justify-center">

            <div className="flex flex-col items-center gap-3">

              <div className="h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />

              <p className="text-sm text-slate-500">
                Loading professors...
              </p>

            </div>

          </CardContent>
        </Card>
      ) : filteredProfessors.length > 0 ? (

        /* =====================================================
           PROFESSORS GRID
        ===================================================== */
        <div className="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">

          {filteredProfessors.map((professor) => {

            const fullName =
              `${professor.first_name || ''} ${
                professor.last_name || ''
              }`.trim() || 'Unknown Professor';

            return (
              <Card
                key={professor.id}
                className="
                  group
                  overflow-hidden
                  transition-all
                  duration-200
                  hover:-translate-y-0.5
                  hover:shadow-md
                "
              >

                <CardContent className="p-5">

                  {/* Professor identity */}
                  <div className="flex items-start gap-4">

                    <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-blue-600 to-cyan-600 text-sm font-bold text-white shadow-sm">
                      {getInitials(professor)}
                    </div>

                    <div className="min-w-0 flex-1">

                      <h3 className="truncate text-base font-semibold text-slate-900">
                        {fullName}
                      </h3>

                      <p className="mt-0.5 text-xs text-slate-500">
                        Professor
                      </p>

                    </div>

                    {/* Status */}
                    <span
                      className={`
                        inline-flex
                        shrink-0
                        items-center
                        gap-1.5
                        rounded-full
                        px-2.5
                        py-1
                        text-[10px]
                        font-semibold
                        uppercase
                        tracking-wide
                        ${
                          professor.active
                            ? 'bg-emerald-50 text-emerald-700'
                            : 'bg-slate-100 text-slate-500'
                        }
                      `}
                    >
                      <span
                        className={`
                          h-1.5
                          w-1.5
                          rounded-full
                          ${
                            professor.active
                              ? 'bg-emerald-500'
                              : 'bg-slate-400'
                          }
                        `}
                      />

                      {professor.active
                        ? 'Active'
                        : 'Inactive'}
                    </span>

                  </div>

                  {/* Information */}
                  <div className="mt-5 space-y-2.5">

                    {/* Email */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">

                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                        <Mail className="h-4 w-4" />
                      </div>

                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Email
                        </p>

                        <p className="truncate text-sm text-slate-700">
                          {professor.email ||
                            'No email available'}
                        </p>
                      </div>

                    </div>

                    {/* Department */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">

                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-cyan-50 text-cyan-600">
                        <Building2 className="h-4 w-4" />
                      </div>

                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Department
                        </p>

                        <p className="truncate text-sm text-slate-700">
                          {professor.department ||
                            'No department'}
                        </p>
                      </div>

                    </div>

                    {/* Grade */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">

                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
                        <User className="h-4 w-4" />
                      </div>

                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Grade
                        </p>

                        <p className="truncate text-sm text-slate-700">
                          {professor.grade ||
                            'No grade'}
                        </p>
                      </div>

                    </div>

                  </div>

                  {/* Footer */}
                  <div className="mt-5 flex items-center justify-end border-t border-slate-100 pt-4">

                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="group/button"
                    >
                      <span>View Details</span>

                      <ArrowUpRight
                        className="
                          h-4
                          w-4
                          transition-transform
                          duration-200
                          group-hover/button:-translate-y-0.5
                          group-hover/button:translate-x-0.5
                        "
                      />
                    </Button>

                  </div>

                </CardContent>

              </Card>
            );
          })}

        </div>

      ) : (

        /* =====================================================
           EMPTY STATE
        ===================================================== */
        <Card>
          <CardContent className="flex min-h-[360px] flex-col items-center justify-center text-center">

            <div className="mb-5 flex h-16 w-16 items-center justify-center rounded-2xl bg-blue-50 text-blue-500">
              <User className="h-7 w-7" />
            </div>

            <h3 className="text-base font-semibold text-slate-900">
              No professors found
            </h3>

            <p className="mt-1 max-w-sm text-sm text-slate-500">
              {searchTerm
                ? 'No professor matches your search criteria.'
                : 'There are currently no professors registered in the system.'}
            </p>

            {searchTerm && (
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="mt-5"
                onClick={() => setSearchTerm('')}
              >
                Clear Search
              </Button>
            )}

          </CardContent>
        </Card>
      )}

    </div>
  );
};