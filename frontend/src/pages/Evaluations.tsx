import React, { useEffect, useState } from 'react';

import {
  evaluationsApi,
  Evaluation,
} from '../api/evaluations';

import {
  coursesApi,
  Course,
} from '../api/courses';

import {
  professorsApi,
  Professor,
} from '../api/professors';

import {
  Card,
  CardContent,
  CardHeader,
} from '../components/ui/Card';

import { Button } from '../components/ui/Button';

import {
  ClipboardList,
  Star,
  CalendarDays,
  User as UserIcon,
  Plus,
  ArrowUpRight,
} from 'lucide-react';

export const Evaluations: React.FC = () => {
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [courses, setCourses] = useState<Course[]>([]);
  const [professors, setProfessors] = useState<Professor[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [evals, coursesData, profsData] = await Promise.all([
        evaluationsApi.getAll(),
        coursesApi.getAll(),
        professorsApi.getActive(),
      ]);

      setEvaluations(evals);
      setCourses(coursesData);
      setProfessors(profsData);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const getProfessorName = (profId: number) => {
    const professor = professors.find(
      (prof) => prof.id === profId
    );

    return professor
      ? `${professor.first_name} ${professor.last_name}`
      : 'Unknown professor';
  };

  const getCourseName = (courseId: number) => {
    const course = courses.find(
      (course) => course.id === courseId
    );

    return course?.name || 'Unknown course';
  };

  const calculateAverageScore = (
    scores: Record<number, number>
  ) => {
    const values = Object.values(scores);

    if (values.length === 0) {
      return '0.0';
    }

    const average =
      values.reduce((total, score) => total + score, 0) /
      values.length;

    return average.toFixed(1);
  };

  return (
    <div className="space-y-7">

      {/* Header */}
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <div className="mb-1 flex items-center gap-2">
            <span className="h-1.5 w-1.5 rounded-full bg-blue-600" />

            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-blue-600">
              Academic Evaluation
            </p>
          </div>

          <h2 className="text-2xl font-bold tracking-tight text-slate-900">
            Evaluations
          </h2>

          <p className="mt-1 text-sm text-slate-500">
            Manage and review professor evaluations.
          </p>
        </div>

        <Button
          type="button"
          className="w-full sm:w-auto"
        >
          <Plus className="h-4 w-4" />
          <span>New Evaluation</span>
        </Button>
      </div>

      {/* Summary */}
      {!isLoading && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Total Evaluations
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {evaluations.length}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                  <ClipboardList className="h-5 w-5" />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Professors Evaluated
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {new Set(
                      evaluations.map(
                        (evaluation) =>
                          evaluation.professor_id
                      )
                    ).size}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-cyan-50 text-cyan-600">
                  <UserIcon className="h-5 w-5" />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-5">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.12em] text-slate-500">
                    Courses Evaluated
                  </p>

                  <p className="mt-2 text-2xl font-bold tracking-tight text-slate-900">
                    {new Set(
                      evaluations.map(
                        (evaluation) =>
                          evaluation.course_id
                      )
                    ).size}
                  </p>
                </div>

                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                  <ClipboardList className="h-5 w-5" />
                </div>
              </div>
            </CardContent>
          </Card>

        </div>
      )}

      {/* Loading */}
      {isLoading ? (
        <Card>
          <CardContent className="flex min-h-[300px] items-center justify-center">
            <div className="flex flex-col items-center gap-3">
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-blue-600" />

              <p className="text-sm text-slate-500">
                Loading evaluations...
              </p>
            </div>
          </CardContent>
        </Card>
      ) : evaluations.length > 0 ? (

        /* Evaluations */
        <div className="grid grid-cols-1 gap-5 xl:grid-cols-2">

          {evaluations.map((evaluation) => {
            const averageScore =
              calculateAverageScore(
                evaluation.criteria_scores
              );

            return (
              <Card
                key={evaluation.id}
                className="group overflow-hidden transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md"
              >

                {/* Card Header */}
                <CardHeader>
                  <div className="flex min-w-0 items-center gap-3">

                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                      <ClipboardList className="h-5 w-5" />
                    </div>

                    <div className="min-w-0">
                      <h3 className="truncate text-sm font-semibold text-slate-900">
                        Evaluation #{evaluation.id}
                      </h3>

                      <p className="mt-0.5 text-xs text-slate-500">
                        Professor evaluation
                      </p>
                    </div>
                  </div>

                  {/* Score */}
                  <div className="flex shrink-0 items-center gap-1.5 rounded-lg bg-amber-50 px-2.5 py-1.5 text-amber-600">
                    <Star className="h-4 w-4 fill-current" />

                    <span className="text-sm font-bold">
                      {averageScore}
                    </span>

                    <span className="text-xs text-amber-500">
                      /5
                    </span>
                  </div>
                </CardHeader>

                {/* Card Content */}
                <CardContent>

                  <div className="space-y-3">

                    {/* Professor */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
                        <UserIcon className="h-4 w-4" />
                      </div>

                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Professor
                        </p>

                        <p className="truncate text-sm font-medium text-slate-800">
                          {getProfessorName(
                            evaluation.professor_id
                          )}
                        </p>
                      </div>
                    </div>

                    {/* Course */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-cyan-50 text-cyan-600">
                        <ClipboardList className="h-4 w-4" />
                      </div>

                      <div className="min-w-0">
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Course
                        </p>

                        <p className="truncate text-sm font-medium text-slate-800">
                          {getCourseName(
                            evaluation.course_id
                          )}
                        </p>
                      </div>
                    </div>

                    {/* Date */}
                    <div className="flex items-center gap-3 rounded-xl border border-slate-100 bg-slate-50/70 px-3 py-2.5">
                      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
                        <CalendarDays className="h-4 w-4" />
                      </div>

                      <div>
                        <p className="text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Date
                        </p>

                        <p className="text-sm font-medium text-slate-800">
                          {new Date(
                            evaluation.created_at
                          ).toLocaleDateString('en-US', {
                            day: '2-digit',
                            month: 'short',
                            year: 'numeric',
                          })}
                        </p>
                      </div>
                    </div>

                    {/* Comments */}
                    {evaluation.comments && (
                      <div className="rounded-xl border border-slate-200 bg-white p-3">
                        <p className="mb-1.5 text-[10px] font-semibold uppercase tracking-[0.12em] text-slate-400">
                          Comments
                        </p>

                        <p className="break-words text-sm leading-6 text-slate-600">
                          {evaluation.comments}
                        </p>
                      </div>
                    )}

                  </div>

                  {/* Footer */}
                  <div className="mt-5 flex justify-end border-t border-slate-100 pt-4">

                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="group/button"
                    >
                      <span>View Details</span>

                      <ArrowUpRight className="h-4 w-4 transition-transform duration-200 group-hover/button:translate-x-0.5 group-hover/button:-translate-y-0.5" />
                    </Button>

                  </div>

                </CardContent>
              </Card>
            );
          })}

        </div>

      ) : (

        /* Empty state */
        <Card>
          <CardContent className="flex min-h-[360px] flex-col items-center justify-center text-center">

            <div className="mb-5 flex h-16 w-16 items-center justify-center rounded-2xl bg-blue-50 text-blue-500">
              <ClipboardList className="h-7 w-7" />
            </div>

            <h3 className="text-base font-semibold text-slate-900">
              No evaluations found
            </h3>

            <p className="mt-1 max-w-sm text-sm text-slate-500">
              There are currently no professor evaluations
              available in the system.
            </p>

            <Button
              type="button"
              className="mt-5"
            >
              <Plus className="h-4 w-4" />
              <span>Create First Evaluation</span>
            </Button>

          </CardContent>
        </Card>
      )}

    </div>
  );
};