import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({
  children,
  className = '',
}) => {
  return (
    <div
      className={`
        rounded-2xl
        border border-slate-200
        bg-white
        shadow-sm
        shadow-slate-200/60
        transition-all
        duration-200
        ${className}
      `}
    >
      {children}
    </div>
  );
};

interface CardHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export const CardHeader: React.FC<CardHeaderProps> = ({
  children,
  className = '',
}) => {
  return (
    <div
      className={`
        flex items-center justify-between
        border-b border-slate-100
        px-5 py-4
        sm:px-6
        ${className}
      `}
    >
      {children}
    </div>
  );
};

interface CardContentProps {
  children: React.ReactNode;
  className?: string;
}

export const CardContent: React.FC<CardContentProps> = ({
  children,
  className = '',
}) => {
  return (
    <div
      className={`
        px-5 py-5
        sm:px-6
        ${className}
      `}
    >
      {children}
    </div>
  );
};
