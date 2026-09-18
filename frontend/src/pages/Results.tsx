import React, { useEffect, useState } from 'react';

import {
  professorsApi,
  Professor,
} from '../api/professors';

import {
  resultsApi,
  ProfessorResult,
} from '../api/results';

import {
  Card,
  CardContent,
  CardHeader,
} from '../components/ui/Card';

import {
  BarChart3,
  TrendingUp,
  Award,
  Star,
  Users,
  Target,
  ChevronRight,
} from 'lucide-react';

export const Results: React.FC = () => {
  const [professors, setProfessors] = useState<Professor[]>([]);
  const [selectedProfessor, setSelectedProfessor] =
    useState<Professor | null>(null);

  const [results, setResults] =
    useState<ProfessorResult | null>(null);

  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingResults, setIsLoadingResults] =
    useState(false);

  useEffect(() => {
    loadProfessors();
  }, []);

  const loadProfessors = async () => {
    try {
      const data = await professorsApi.getActive();

      setProfessors(data);
    } catch (error) {
      console.error(
        'Failed to load professors:',
        error
      );
    } finally {
      setIsLoading(false);
    }
  };

  const loadResults = async (professor: Professor) => {
    setIsLoadingResults(true);

    try {
      const data =
        await resultsApi.getProfessorResults(
          professor.id
        );

      setResults(data);
      setSelectedProfessor(professor);
    } catch (error) {
      console.error(
        'Failed to load results:',
        error
      );

      setResults(null);
      setSelectedProfessor(professor);
    } finally {
      setIsLoadingResults(false);
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 4.5) {
      return 'text-emerald-600';
    }

    if (score >= 3.5) {
      return 'text-blue-600';
    }

    if (score >= 2.5) {
      return 'text-amber-600';
    }

    return 'text-red-600';
  };

  const getScoreBackground = (score: number) => {
    if (score >= 4.5) {
      return 'bg-emerald-500';
    }

    if (score >= 3.5) {
      return 'bg-blue-600';
    }

    if (score >= 2.5) {
      return 'bg-amber-500';
    }

    return 'bg-red-500';
  };

  const getScoreLightBackground = (score: number) => {
    if (score >= 4.5) {
      return 'bg-emerald-50';
    }

    if (score >= 3.5) {
      return 'bg-blue-50';
    }

    if (score >= 2.5) {
      return 'bg-amber-50';
    }

    return 'bg-red-50';
  };

  const getPerformanceLabel = (score: number) => {
    if (score >= 4.5) {
      return 'Outstanding Performance';
    }

    if (score >= 4.0) {
      return 'Excellent Performance';
    }

    if (score >= 3.0) {
      return 'Good Performance';
    }

    if (score >= 2.5) {
      return 'Moderate Performance';
    }

    return 'Needs Improvement';
  };

  return (
    <div className="space-y-7">

      {/* =====================================================
          HEADER
      ===================================================== */}
      <div>
        <div className="mb-1 flex items-center gap-2">
          <span className="h-1.5 w-1.5 rounded-full bg-blue-600" />

          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-blue-600">
            Analytics
          </p>
        </div>

        <h2 className="text-2xl font-bold tracking-tight text-slate-900">
          Results & Analytics
        </h2>

        <p className="mt-1 text-sm text-slate-500">
          View professor evaluation results and performance
          statistics.
        </p>
      </div>

      {/* =====================================================
          MAIN LAYOUT
      ===================================================== */}
      <div className="grid grid-cols-1 gap-5 xl:grid-cols-[300px_1fr]">

        {/* ===================================================
            PROFESSOR SELECTOR
        =================================================== */}
        <Card className="h-fit">

          <CardHeader>
            <div>
              <h3 className="text-sm font-semibold text-slate-900">
                Select Professor
              </h3>

              <p className="mt-0.5 text-xs text-slate-500">
                Choose a professor to view results
              </p>
            </div>

            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
              <Users className="h-4 w-4" />
            </div>
          </CardHeader>

          <CardContent className="p-3">

            {isLoading ? (
              <div className="flex flex-col items-center justify-center py-10">

                <div className="h-7 w-7 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />

                <p className="mt-3 text-xs text-slate-500">
                  Loading professors...
                </p>

              </div>
            ) : professors.length === 0 ? (
              <div className="py-10 text-center">

                <Users className="mx-auto h-8 w-8 text-slate-300" />

                <p className="mt-3 text-sm font-medium text-slate-700">
                  No professors found
                </p>

                <p className="mt-1 text-xs text-slate-500">
                  There are no active professors available.
                </p>

              </div>
            ) : (
              <div className="max-h-[560px] space-y-1.5 overflow-y-auto pr-1">

                {professors.map((professor) => {

                  const isSelected =
                    selectedProfessor?.id === professor.id;

                  return (
                    <button
                      key={professor.id}
                      type="button"
                      onClick={() =>
                        loadResults(professor)
                      }
                      className={`
                        group
                        flex
                        w-full
                        items-center
                        gap-3
                        rounded-xl
                        border
                        p-3
                        text-left
                        transition-all
                        duration-200
                        ${
                          isSelected
                            ? 'border-blue-200 bg-blue-50'
                            : 'border-transparent bg-slate-50 hover:border-slate-200 hover:bg-white'
                        }
                      `}
                    >

                      {/* Avatar */}
                      <div
                        className={`
                          flex
                          h-9
                          w-9
                          shrink-0
                          items-center
                          justify-center
                          rounded-lg
                          text-xs
                          font-bold
                          ${
                            isSelected
                              ? 'bg-blue-600 text-white'
                              : 'bg-slate-200 text-slate-600'
                          }
                        `}
                      >
                        {professor.first_name?.charAt(0)}
                        {professor.last_name?.charAt(0)}
                      </div>

                      {/* Name */}
                      <div className="min-w-0 flex-1">

                        <p
                          className={`
                            truncate
                            text-sm
                            font-semibold
                            ${
                              isSelected
                                ? 'text-blue-700'
                                : 'text-slate-800'
                            }
                          `}
                        >
                          {professor.first_name}{' '}
                          {professor.last_name}
                        </p>

                        <p
                          className={`
                            mt-0.5
                            truncate
                            text-xs
                            ${
                              isSelected
                                ? 'text-blue-600/70'
                                : 'text-slate-500'
                            }
                          `}
                        >
                          {professor.department ||
                            'No department'}
                        </p>

                      </div>

                      <ChevronRight
                        className={`
                          h-4
                          w-4
                          shrink-0
                          transition-transform
                          ${
                            isSelected
                              ? 'translate-x-0 text-blue-500'
                              : '-translate-x-1 text-slate-300 opacity-0 group-hover:translate-x-0 group-hover:opacity-100'
                          }
                        `}
                      />

                    </button>
                  );
                })}

              </div>
            )}

          </CardContent>
        </Card>

        {/* ===================================================
            RESULTS
        =================================================== */}
        <div className="min-w-0">

          {isLoadingResults ? (

            <Card>
              <CardContent className="flex min-h-[400px] items-center justify-center">

                <div className="flex flex-col items-center gap-3">

                  <div className="h-9 w-9 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />

                  <p className="text-sm text-slate-500">
                    Loading results...
                  </p>

                </div>

              </CardContent>
            </Card>

          ) : results ? (

            <div className="space-y-5">

              {/* =================================================
                  OVERVIEW
              ================================================= */}
              <Card className="overflow-hidden">

                <CardHeader>

                  <div className="min-w-0">

                    <p className="mb-1 text-[10px] font-semibold uppercase tracking-[0.15em] text-blue-600">
                      Professor Results
                    </p>

                    <h3 className="truncate text-lg font-bold text-slate-900">
                      {results.professor_name}
                    </h3>

                  </div>

                  <div
                    className={`
                      flex
                      items-center
                      gap-2
                      rounded-xl
                      px-3
                      py-2
                      ${getScoreLightBackground(
                        results.average_score
                      )}
                    `}
                  >
                    <Star
                      className={`
                        h-5
                        w-5
                        fill-current
                        ${getScoreColor(
                          results.average_score
                        )}
                      `}
                    />

                    <span
                      className={`
                        text-2xl
                        font-bold
                        ${getScoreColor(
                          results.average_score
                        )}
                      `}
                    >
                      {results.average_score.toFixed(1)}
                    </span>

                    <span className="text-xs text-slate-500">
                      /5
                    </span>
                  </div>

                </CardHeader>

                <CardContent>

                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">

                    {/* Total evaluations */}
                    <div className="rounded-xl border border-slate-100 bg-slate-50/70 p-4">

                      <div className="flex items-center justify-between">

                        <div>
                          <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                            Evaluations
                          </p>

                          <p className="mt-2 text-2xl font-bold text-slate-900">
                            {results.total_evaluations}
                          </p>
                        </div>

                        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                          <Users className="h-5 w-5" />
                        </div>

                      </div>

                    </div>

                    {/* Average score */}
                    <div className="rounded-xl border border-slate-100 bg-slate-50/70 p-4">

                      <div className="flex items-center justify-between">

                        <div>
                          <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                            Average Score
                          </p>

                          <p
                            className={`
                              mt-2
                              text-2xl
                              font-bold
                              ${getScoreColor(
                                results.average_score
                              )}
                            `}
                          >
                            {results.average_score.toFixed(1)}
                          </p>
                        </div>

                        <div
                          className={`
                            flex
                            h-10
                            w-10
                            items-center
                            justify-center
                            rounded-xl
                            ${getScoreLightBackground(
                              results.average_score
                            )}
                          `}
                        >
                          <TrendingUp className="h-5 w-5" />
                        </div>

                      </div>

                    </div>

                    {/* Performance */}
                    <div className="rounded-xl border border-slate-100 bg-slate-50/70 p-4">

                      <div className="flex items-center justify-between">

                        <div className="min-w-0">

                          <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                            Performance
                          </p>

                          <p className="mt-2 truncate text-sm font-bold text-slate-900">
                            {getPerformanceLabel(
                              results.average_score
                            )}
                          </p>

                        </div>

                        <div
                          className={`
                            flex
                            h-10
                            w-10
                            shrink-0
                            items-center
                            justify-center
                            rounded-xl
                            ${getScoreLightBackground(
                              results.average_score
                            )}
                            ${getScoreColor(
                              results.average_score
                            )}
                          `}
                        >
                          <Award className="h-5 w-5" />
                        </div>

                      </div>

                    </div>

                  </div>

                </CardContent>
              </Card>

              {/* =================================================
                  CRITERIA BREAKDOWN
              ================================================= */}
              <Card>

                <CardHeader>

                  <div>
                    <h3 className="text-sm font-semibold text-slate-900">
                      Criteria Breakdown
                    </h3>

                    <p className="mt-0.5 text-xs text-slate-500">
                      Average score for each evaluation criterion
                    </p>
                  </div>

                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-cyan-50 text-cyan-600">
                    <Target className="h-4 w-4" />
                  </div>

                </CardHeader>

                <CardContent>

                  <div className="space-y-5">

                    {Object.entries(
                      results.criteria_averages
                    ).map(([criterion, score]) => {

                      const percentage =
                        Math.min(
                          Math.max(
                            (score / 5) * 100,
                            0
                          ),
                          100
                        );

                      return (
                        <div key={criterion}>

                          <div className="mb-2 flex items-center justify-between gap-4">

                            <span className="min-w-0 truncate text-sm font-medium text-slate-700">
                              {criterion}
                            </span>

                            <span
                              className={`
                                shrink-0
                                text-sm
                                font-bold
                                ${getScoreColor(score)}
                              `}
                            >
                              {score.toFixed(1)}
                              <span className="ml-0.5 text-xs font-medium text-slate-400">
                                /5
                              </span>
                            </span>

                          </div>

                          <div className="h-2.5 w-full overflow-hidden rounded-full bg-slate-100">

                            <div
                              className={`
                                h-full
                                rounded-full
                                transition-all
                                duration-500
                                ${getScoreBackground(score)}
                              `}
                              style={{
                                width: `${percentage}%`,
                              }}
                            />

                          </div>

                        </div>
                      );
                    })}

                  </div>

                </CardContent>
              </Card>

              {/* =================================================
                  PERFORMANCE SUMMARY
              ================================================= */}
              <Card>

                <CardHeader>

                  <div className="flex items-center gap-3">

                    <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                      <TrendingUp className="h-5 w-5" />
                    </div>

                    <div>
                      <h3 className="text-sm font-semibold text-slate-900">
                        Performance Summary
                      </h3>

                      <p className="mt-0.5 text-xs text-slate-500">
                        Overall evaluation assessment
                      </p>
                    </div>

                  </div>

                </CardHeader>

                <CardContent>

                  <div
                    className={`
                      flex
                      items-center
                      gap-4
                      rounded-xl
                      border
                      p-4
                      ${getScoreLightBackground(
                        results.average_score
                      )}
                    `}
                  >

                    <div
                      className={`
                        flex
                        h-12
                        w-12
                        shrink-0
                        items-center
                        justify-center
                        rounded-xl
                        bg-white
                        shadow-sm
                        ${getScoreColor(
                          results.average_score
                        )}
                      `}
                    >
                      <Award className="h-6 w-6" />
                    </div>

                    <div className="min-w-0">

                      <p
                        className={`
                          text-sm
                          font-bold
                          ${getScoreColor(
                            results.average_score
                          )}
                        `}
                      >
                        {getPerformanceLabel(
                          results.average_score
                        )}
                      </p>

                      <p className="mt-1 text-xs text-slate-500">
                        Based on{' '}
                        <span className="font-semibold text-slate-700">
                          {results.total_evaluations}
                        </span>{' '}
                        evaluation
                        {results.total_evaluations !== 1
                          ? 's'
                          : ''}
                        .
                      </p>

                    </div>

                  </div>

                </CardContent>
              </Card>

            </div>

          ) : (

            /* =================================================
               EMPTY STATE
            ================================================= */
            <Card>
              <CardContent className="flex min-h-[500px] flex-col items-center justify-center text-center">

                <div className="mb-5 flex h-16 w-16 items-center justify-center rounded-2xl bg-blue-50 text-blue-500">
                  <BarChart3 className="h-7 w-7" />
                </div>

                <h3 className="text-base font-semibold text-slate-900">
                  Select a professor
                </h3>

                <p className="mt-1 max-w-sm text-sm text-slate-500">
                  Select a professor from the list to view
                  their evaluation results and performance
                  statistics.
                </p>

              </CardContent>
            </Card>
          )}

        </div>

      </div>
    </div>
  );
};