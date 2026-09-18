import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import {
  Users,
  ClipboardList,
  TrendingUp,
  Award,
  Plus,
  BarChart3,
  User as UserIcon,
  ArrowUpRight,
  CalendarDays,
} from 'lucide-react';

import { Card, CardContent, CardHeader } from '../components/ui/Card';
import { professorsApi, Professor } from '../api/professors';
import { evaluationsApi, Evaluation } from '../api/evaluations';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  const [professors, setProfessors] = useState<Professor[]>([]);
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [professorsData, evaluationsData] = await Promise.all([
        professorsApi.getActive(),
        evaluationsApi.getAll(),
      ]);

      setProfessors(professorsData);
      setEvaluations(evaluationsData);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const getTopRatedProfessor = () => {
    if (professors.length === 0) {
      return {
        name: 'N/A',
        score: 0,
      };
    }

    return {
      name: `${professors[0].first_name} ${professors[0].last_name}`,
      score: 4.5,
    };
  };

  const topProfessor = getTopRatedProfessor();

  const stats = [
    {
      label: 'Total Professors',
      value: professors.length.toString(),
      icon: <Users className="h-5 w-5" />,
      iconBg: 'bg-blue-50',
      iconColor: 'text-blue-600',
    },
    {
      label: 'Active Evaluations',
      value: evaluations.length.toString(),
      icon: <ClipboardList className="h-5 w-5" />,
      iconBg: 'bg-cyan-50',
      iconColor: 'text-cyan-600',
    },
    {
      label: 'Average Score',
      value: '4.2/5',
      icon: <TrendingUp className="h-5 w-5" />,
      iconBg: 'bg-emerald-50',
      iconColor: 'text-emerald-600',
    },
    {
      label: 'Top Rated',
      value: topProfessor.name,
      icon: <Award className="h-5 w-5" />,
      iconBg: 'bg-amber-50',
      iconColor: 'text-amber-600',
    },
  ];

  if (isLoading) {
    return (
      <div className="flex min-h-[300px] items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />
          <p className="text-sm text-slate-500">
            Loading dashboard...
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-7">

      {/* =====================================================
          WELCOME SECTION
      ====================================================== */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <p className="mb-1 text-xs font-semibold uppercase tracking-[0.18em] text-blue-600">
            Overview
          </p>

          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Dashboard
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Overview of the professor evaluation system.
          </p>
        </div>

        <div className="flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-500 shadow-sm">
          <CalendarDays className="h-4 w-4 text-blue-600" />

          <span>
            {new Date().toLocaleDateString('en-US', {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
              year: 'numeric',
            })}
          </span>
        </div>
      </div>

      {/* =====================================================
          STATISTICS
      ====================================================== */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">

        {stats.map((stat) => (
          <Card
            key={stat.label}
            className="
              overflow-hidden
              transition-all
              duration-200
              hover:-translate-y-0.5
              hover:shadow-md
            "
          >
            <CardContent className="p-5">

              <div className="flex items-start justify-between gap-4">

                <div className="min-w-0">
                  <p className="mb-2 text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    {stat.label}
                  </p>

                  <p
                    className={`
                      truncate
                      text-2xl
                      font-bold
                      tracking-tight
                      text-slate-900

                      ${
                        stat.label === 'Top Rated'
                          ? 'text-base sm:text-lg'
                          : ''
                      }
                    `}
                  >
                    {stat.value}
                  </p>
                </div>

                <div
                  className={`
                    flex
                    h-11
                    w-11
                    shrink-0
                    items-center
                    justify-center
                    rounded-xl
                    ${stat.iconBg}
                    ${stat.iconColor}
                  `}
                >
                  {stat.icon}
                </div>

              </div>
            </CardContent>
          </Card>
        ))}

      </div>

      {/* =====================================================
          MAIN CONTENT
      ====================================================== */}
      <div className="grid grid-cols-1 gap-5 xl:grid-cols-5">

        {/* Recent activity */}
        <Card className="xl:col-span-3">

          <CardHeader>
            <div>
              <h3 className="text-base font-semibold text-slate-900">
                Recent Activity
              </h3>

              <p className="mt-0.5 text-xs text-slate-500">
                Latest evaluations submitted
              </p>
            </div>

            <button
              type="button"
              onClick={() => navigate('/evaluations')}
              className="
                flex
                items-center
                gap-1
                text-xs
                font-semibold
                text-blue-600
                transition-colors
                hover:text-blue-700
              "
            >
              View all
              <ArrowUpRight className="h-3.5 w-3.5" />
            </button>
          </CardHeader>

          <CardContent>

            {evaluations.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">

                <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-slate-100">
                  <ClipboardList className="h-6 w-6 text-slate-400" />
                </div>

                <p className="text-sm font-medium text-slate-700">
                  No recent evaluations
                </p>

                <p className="mt-1 text-xs text-slate-500">
                  Evaluations will appear here once submitted.
                </p>

              </div>
            ) : (
              <div className="space-y-2">

                {evaluations.slice(0, 5).map((evaluation) => (
                  <div
                    key={evaluation.id}
                    className="
                      group
                      flex
                      items-center
                      gap-3
                      rounded-xl
                      border
                      border-slate-100
                      bg-slate-50/70
                      p-3
                      transition-all
                      duration-200
                      hover:border-blue-100
                      hover:bg-blue-50/40
                    "
                  >

                    <div
                      className="
                        flex
                        h-9
                        w-9
                        shrink-0
                        items-center
                        justify-center
                        rounded-lg
                        bg-blue-600
                        text-xs
                        font-bold
                        text-white
                      "
                    >
                      #{evaluation.id}
                    </div>

                    <div className="min-w-0 flex-1">

                      <p className="truncate text-sm font-semibold text-slate-800">
                        Evaluation #{evaluation.id}
                      </p>

                      <div className="mt-1 flex items-center gap-1.5 text-xs text-slate-500">
                        <CalendarDays className="h-3.5 w-3.5" />

                        {new Date(
                          evaluation.created_at
                        ).toLocaleDateString()}
                      </div>

                    </div>

                    <ArrowUpRight
                      className="
                        h-4
                        w-4
                        text-slate-300
                        transition-colors
                        group-hover:text-blue-500
                      "
                    />

                  </div>
                ))}

              </div>
            )}

          </CardContent>
        </Card>

        {/* =====================================================
            QUICK ACTIONS
        ====================================================== */}
        <Card className="xl:col-span-2">

          <CardHeader>
            <div>
              <h3 className="text-base font-semibold text-slate-900">
                Quick Actions
              </h3>

              <p className="mt-0.5 text-xs text-slate-500">
                Frequently used actions
              </p>
            </div>
          </CardHeader>

          <CardContent>

            <div className="space-y-2.5">

              {/* New Evaluation */}
              <button
                type="button"
                onClick={() => navigate('/evaluations')}
                className="
                  group
                  flex
                  w-full
                  items-center
                  gap-3
                  rounded-xl
                  border
                  border-slate-200
                  bg-white
                  p-3
                  text-left
                  transition-all
                  duration-200
                  hover:border-blue-200
                  hover:bg-blue-50/50
                "
              >
                <div
                  className="
                    flex
                    h-10
                    w-10
                    shrink-0
                    items-center
                    justify-center
                    rounded-xl
                    bg-blue-50
                    text-blue-600
                    transition-colors
                    group-hover:bg-blue-100
                  "
                >
                  <Plus className="h-5 w-5" />
                </div>

                <div className="min-w-0 flex-1">
                  <p className="text-sm font-semibold text-slate-800">
                    New Evaluation
                  </p>

                  <p className="mt-0.5 text-xs text-slate-500">
                    Submit a professor evaluation
                  </p>
                </div>

                <ArrowUpRight className="h-4 w-4 text-slate-300 group-hover:text-blue-500" />
              </button>

              {/* Results */}
              <button
                type="button"
                onClick={() => navigate('/results')}
                className="
                  group
                  flex
                  w-full
                  items-center
                  gap-3
                  rounded-xl
                  border
                  border-slate-200
                  bg-white
                  p-3
                  text-left
                  transition-all
                  duration-200
                  hover:border-cyan-200
                  hover:bg-cyan-50/40
                "
              >
                <div
                  className="
                    flex
                    h-10
                    w-10
                    shrink-0
                    items-center
                    justify-center
                    rounded-xl
                    bg-cyan-50
                    text-cyan-600
                    transition-colors
                    group-hover:bg-cyan-100
                  "
                >
                  <BarChart3 className="h-5 w-5" />
                </div>

                <div className="min-w-0 flex-1">
                  <p className="text-sm font-semibold text-slate-800">
                    View Results
                  </p>

                  <p className="mt-0.5 text-xs text-slate-500">
                    See evaluation statistics
                  </p>
                </div>

                <ArrowUpRight className="h-4 w-4 text-slate-300 group-hover:text-cyan-500" />
              </button>

              {/* Professors */}
              <button
                type="button"
                onClick={() => navigate('/professors')}
                className="
                  group
                  flex
                  w-full
                  items-center
                  gap-3
                  rounded-xl
                  border
                  border-slate-200
                  bg-white
                  p-3
                  text-left
                  transition-all
                  duration-200
                  hover:border-emerald-200
                  hover:bg-emerald-50/40
                "
              >
                <div
                  className="
                    flex
                    h-10
                    w-10
                    shrink-0
                    items-center
                    justify-center
                    rounded-xl
                    bg-emerald-50
                    text-emerald-600
                    transition-colors
                    group-hover:bg-emerald-100
                  "
                >
                  <UserIcon className="h-5 w-5" />
                </div>

                <div className="min-w-0 flex-1">
                  <p className="text-sm font-semibold text-slate-800">
                    Manage Professors
                  </p>

                  <p className="mt-0.5 text-xs text-slate-500">
                    View and manage professors
                  </p>
                </div>

                <ArrowUpRight className="h-4 w-4 text-slate-300 group-hover:text-emerald-500" />
              </button>

            </div>

          </CardContent>
        </Card>

      </div>
    </div>
  );
};
