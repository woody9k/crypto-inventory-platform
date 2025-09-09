import React from 'react';

interface RiskBadgeProps {
  riskLevel: string;
  riskScore?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const RiskBadge: React.FC<RiskBadgeProps> = ({ 
  riskLevel, 
  riskScore, 
  size = 'md', 
  className = '' 
}) => {
  const getRiskConfig = (level: string) => {
    switch (level.toLowerCase()) {
      case 'high':
        return {
          bgColor: 'bg-red-100 dark:bg-red-900',
          textColor: 'text-red-800 dark:text-red-200',
          borderColor: 'border-red-200 dark:border-red-700',
          icon: '🔴',
          label: 'High Risk'
        };
      case 'medium':
        return {
          bgColor: 'bg-yellow-100 dark:bg-yellow-900',
          textColor: 'text-yellow-800 dark:text-yellow-200',
          borderColor: 'border-yellow-200 dark:border-yellow-700',
          icon: '🟡',
          label: 'Medium Risk'
        };
      case 'low':
        return {
          bgColor: 'bg-green-100 dark:bg-green-900',
          textColor: 'text-green-800 dark:text-green-200',
          borderColor: 'border-green-200 dark:border-green-700',
          icon: '🟢',
          label: 'Low Risk'
        };
      default:
        return {
          bgColor: 'bg-gray-100 dark:bg-gray-800',
          textColor: 'text-gray-800 dark:text-gray-200',
          borderColor: 'border-gray-200 dark:border-gray-600',
          icon: '⚪',
          label: 'Unknown'
        };
    }
  };

  const getSizeClasses = (size: string) => {
    switch (size) {
      case 'sm':
        return 'px-2 py-1 text-xs';
      case 'lg':
        return 'px-4 py-2 text-base';
      default:
        return 'px-3 py-1 text-sm';
    }
  };

  const config = getRiskConfig(riskLevel);
  const sizeClasses = getSizeClasses(size);

  return (
    <span
      className={`
        inline-flex items-center rounded-full border font-medium
        ${config.bgColor}
        ${config.textColor}
        ${config.borderColor}
        ${sizeClasses}
        ${className}
      `}
      title={riskScore ? `Risk Score: ${riskScore}` : config.label}
    >
      <span className="mr-1" role="img" aria-label="risk indicator">
        {config.icon}
      </span>
      {config.label}
      {riskScore !== undefined && (
        <span className="ml-1 font-semibold">
          ({riskScore})
        </span>
      )}
    </span>
  );
};
