import React from 'react';

interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  className = '',
  ...props
}) => {
  return (
    <div className="w-full">

      {/* Label */}
      {label && (
        <label
          htmlFor={props.id}
          className="
            mb-1.5
            block
            text-sm
            font-medium
            text-slate-700
          "
        >
          {label}
        </label>
      )}

      {/* Input */}
      <input
        {...props}
        className={`
          w-full
          rounded-xl
          border
          bg-white
          px-3.5
          py-2.5
          text-sm
          text-slate-900

          placeholder:text-slate-400

          transition-all
          duration-200

          focus:outline-none
          focus:ring-4

          hover:border-slate-300

          disabled:cursor-not-allowed
          disabled:bg-slate-100
          disabled:text-slate-500

          ${
            error
              ? `
                border-red-400
                focus:border-red-500
                focus:ring-red-100
              `
              : `
                border-slate-200
                focus:border-blue-500
                focus:ring-blue-100
              `
          }

          ${className}
        `}
      />

      {/* Error */}
      {error && (
        <p
          className="
            mt-1.5
            text-xs
            font-medium
            text-red-600
          "
        >
          {error}
        </p>
      )}
    </div>
  );
};
