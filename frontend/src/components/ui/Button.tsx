import React from 'react';

interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  className = '',
  children,
  ...props
}) => {
  const baseStyles = `
    inline-flex
    items-center
    justify-center
    gap-2
    rounded-xl
    font-medium
    transition-all
    duration-200
    focus:outline-none
    focus:ring-2
    focus:ring-offset-2
    disabled:pointer-events-none
    disabled:cursor-not-allowed
    disabled:opacity-50
    active:scale-[0.98]
  `;

  const variants = {
    primary: `
      bg-blue-600
      text-white
      shadow-md
      shadow-blue-600/15
      hover:bg-blue-700
      hover:shadow-lg
      hover:shadow-blue-600/20
      focus:ring-blue-500
    `,

    secondary: `
      bg-slate-800
      text-white
      shadow-sm
      hover:bg-slate-900
      focus:ring-slate-500
    `,

    outline: `
      border
      border-blue-600
      bg-white
      text-blue-600
      hover:bg-blue-50
      focus:ring-blue-500
    `,

    ghost: `
      border
      border-transparent
      bg-transparent
      text-slate-600
      hover:bg-slate-100
      hover:text-blue-600
      focus:ring-slate-400
    `,

    danger: `
      bg-red-600
      text-white
      shadow-md
      shadow-red-600/15
      hover:bg-red-700
      hover:shadow-lg
      hover:shadow-red-600/20
      focus:ring-red-500
    `,
  };

  const sizes = {
    sm: `
      min-h-9
      px-3
      text-sm
    `,

    md: `
      min-h-10
      px-4
      text-sm
    `,

    lg: `
      min-h-12
      px-6
      text-base
    `,
  };

  return (
    <button
      className={`
        ${baseStyles}
        ${variants[variant]}
        ${sizes[size]}
        ${className}
      `}
      {...props}
    >
      {children}
    </button>
  );
};
