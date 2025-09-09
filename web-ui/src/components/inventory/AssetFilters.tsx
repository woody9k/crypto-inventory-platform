import React, { useState } from 'react';
import { MagnifyingGlassIcon, FunnelIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { AssetFilters } from '../../types/inventory';
import { Button } from '../common/Button';
import { Input } from '../common/Input';

interface AssetFiltersProps {
  filters: AssetFilters;
  onFiltersChange: (filters: AssetFilters) => void;
  onClearFilters: () => void;
}

export const AssetFiltersComponent: React.FC<AssetFiltersProps> = ({
  filters,
  onFiltersChange,
  onClearFilters,
}) => {
  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onFiltersChange({ ...filters, search: e.target.value, page: 1 });
  };

  const handleFilterChange = (key: keyof AssetFilters, value: string[]) => {
    onFiltersChange({ ...filters, [key]: value, page: 1 });
  };

  const handleSingleFilterChange = (key: keyof AssetFilters, value: string) => {
    onFiltersChange({ ...filters, [key]: value, page: 1 });
  };

  const assetTypes = ['server', 'endpoint', 'service', 'appliance'];
  const environments = ['production', 'staging', 'development', 'test'];
  const riskLevels = ['high', 'medium', 'low', 'unknown'];
  const protocols = ['TLS', 'SSH', 'IPSec', 'VPN', 'Database', 'API'];
  const sortOptions = [
    { value: 'hostname', label: 'Hostname' },
    { value: 'risk_score', label: 'Risk Score' },
    { value: 'last_seen_at', label: 'Last Seen' },
    { value: 'created_at', label: 'Created' },
  ];

  const activeFilterCount = [
    filters.asset_type?.length || 0,
    filters.environment?.length || 0,
    filters.risk_level?.length || 0,
    filters.protocol?.length || 0,
    filters.business_unit?.length || 0,
  ].reduce((sum, count) => sum + count, 0);

  return (
    <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4 mb-6">
      {/* Search Bar */}
      <div className="flex flex-col sm:flex-row gap-4 mb-4">
        <div className="flex-1 relative">
          <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
          <Input
            type="text"
            placeholder="Search assets by hostname, IP, description..."
            value={filters.search || ''}
            onChange={handleSearchChange}
            className="pl-10"
          />
        </div>
        
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={() => setShowAdvanced(!showAdvanced)}
            className="flex items-center gap-2"
          >
            <FunnelIcon className="h-4 w-4" />
            Filters
            {activeFilterCount > 0 && (
              <span className="bg-primary-600 text-white text-xs rounded-full px-2 py-1 min-w-[20px] h-5 flex items-center justify-center">
                {activeFilterCount}
              </span>
            )}
          </Button>
          
          {activeFilterCount > 0 && (
            <Button
              variant="secondary"
              onClick={onClearFilters}
              className="flex items-center gap-2"
            >
              <XMarkIcon className="h-4 w-4" />
              Clear
            </Button>
          )}
        </div>
      </div>

      {/* Advanced Filters */}
      {showAdvanced && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pt-4 border-t border-gray-200 dark:border-gray-700">
          {/* Asset Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Asset Type
            </label>
            <div className="space-y-2">
              {assetTypes.map((type) => (
                <label key={type} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.asset_type?.includes(type) || false}
                    onChange={(e) => {
                      const current = filters.asset_type || [];
                      const updated = e.target.checked
                        ? [...current, type]
                        : current.filter(t => t !== type);
                      handleFilterChange('asset_type', updated);
                    }}
                    className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                  />
                  <span className="ml-2 text-sm text-gray-700 dark:text-gray-300 capitalize">
                    {type}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Environment */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Environment
            </label>
            <div className="space-y-2">
              {environments.map((env) => (
                <label key={env} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.environment?.includes(env) || false}
                    onChange={(e) => {
                      const current = filters.environment || [];
                      const updated = e.target.checked
                        ? [...current, env]
                        : current.filter(e => e !== env);
                      handleFilterChange('environment', updated);
                    }}
                    className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                  />
                  <span className="ml-2 text-sm text-gray-700 dark:text-gray-300 capitalize">
                    {env}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Risk Level */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Risk Level
            </label>
            <div className="space-y-2">
              {riskLevels.map((level) => (
                <label key={level} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.risk_level?.includes(level) || false}
                    onChange={(e) => {
                      const current = filters.risk_level || [];
                      const updated = e.target.checked
                        ? [...current, level]
                        : current.filter(l => l !== level);
                      handleFilterChange('risk_level', updated);
                    }}
                    className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                  />
                  <span className="ml-2 text-sm text-gray-700 dark:text-gray-300 capitalize">
                    {level}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Protocol */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Protocol
            </label>
            <div className="space-y-2">
              {protocols.map((protocol) => (
                <label key={protocol} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.protocol?.includes(protocol) || false}
                    onChange={(e) => {
                      const current = filters.protocol || [];
                      const updated = e.target.checked
                        ? [...current, protocol]
                        : current.filter(p => p !== protocol);
                      handleFilterChange('protocol', updated);
                    }}
                    className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                  />
                  <span className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                    {protocol}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Sort Options */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Sort By
            </label>
            <select
              value={filters.sort_by || 'hostname'}
              onChange={(e) => handleSingleFilterChange('sort_by', e.target.value)}
              className="block w-full rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100"
            >
              {sortOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
            <div className="mt-2">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={filters.sort_order === 'desc'}
                  onChange={(e) => handleSingleFilterChange('sort_order', e.target.checked ? 'desc' : 'asc')}
                  className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                  Descending
                </span>
              </label>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
