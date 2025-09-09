import React from 'react';
import { Asset } from '../../types/inventory';
import { RiskBadge } from './RiskBadge';
import { 
  ComputerDesktopIcon, 
  ServerIcon, 
  CpuChipIcon,
  WrenchScrewdriverIcon,
  EyeIcon 
} from '@heroicons/react/24/outline';

interface AssetTableProps {
  assets: Asset[];
  loading?: boolean;
  onAssetClick?: (asset: Asset) => void;
}

export const AssetTable: React.FC<AssetTableProps> = ({ 
  assets, 
  loading = false, 
  onAssetClick 
}) => {
  const getAssetIcon = (assetType: string) => {
    switch (assetType.toLowerCase()) {
      case 'server':
        return <ServerIcon className="h-5 w-5 text-gray-500" />;
      case 'endpoint':
        return <ComputerDesktopIcon className="h-5 w-5 text-gray-500" />;
      case 'service':
        return <CpuChipIcon className="h-5 w-5 text-gray-500" />;
      case 'appliance':
        return <WrenchScrewdriverIcon className="h-5 w-5 text-gray-500" />;
      default:
        return <ComputerDesktopIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getEnvironmentBadge = (environment?: string) => {
    if (!environment) return null;
    
    const colors = {
      production: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
      staging: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
      development: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
      test: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    };

    const colorClass = colors[environment as keyof typeof colors] || 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200';

    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorClass}`}>
        {environment}
      </span>
    );
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
        <div className="animate-pulse">
          <div className="px-4 py-5 sm:p-6">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4 mb-4"></div>
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex space-x-4 mb-4">
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/6"></div>
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/6"></div>
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (assets.length === 0) {
    return (
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
        <div className="text-center py-12">
          <ComputerDesktopIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
            No assets found
          </h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Try adjusting your search criteria or filters.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-900">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Asset
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Location
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Environment
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Risk Level
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Last Seen
              </th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Crypto Protocols
              </th>
              <th scope="col" className="relative px-6 py-3">
                <span className="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {assets.map((asset) => (
              <tr 
                key={asset.id}
                className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer"
                onClick={() => onAssetClick?.(asset)}
              >
                {/* Asset Info */}
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      {getAssetIcon(asset.asset_type)}
                    </div>
                    <div className="ml-4">
                      <div className="text-sm font-medium text-gray-900 dark:text-white">
                        {asset.hostname || 'Unknown Host'}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400 capitalize">
                        {asset.asset_type}
                        {asset.operating_system && (
                          <span className="ml-2">• {asset.operating_system}</span>
                        )}
                      </div>
                    </div>
                  </div>
                </td>

                {/* Location */}
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-gray-900 dark:text-white">
                    {asset.ip_address || 'N/A'}
                    {asset.port && (
                      <span className="text-gray-500 dark:text-gray-400">:{asset.port}</span>
                    )}
                  </div>
                  {asset.business_unit && (
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      {asset.business_unit}
                    </div>
                  )}
                </td>

                {/* Environment */}
                <td className="px-6 py-4 whitespace-nowrap">
                  {getEnvironmentBadge(asset.environment)}
                </td>

                {/* Risk Level */}
                <td className="px-6 py-4 whitespace-nowrap">
                  <RiskBadge 
                    riskLevel={asset.risk_level} 
                    riskScore={asset.risk_score}
                    size="sm"
                  />
                </td>

                {/* Last Seen */}
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                  {formatDate(asset.last_seen_at)}
                </td>

                {/* Crypto Protocols */}
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex flex-wrap gap-1">
                    {asset.crypto_implementations?.slice(0, 3).map((crypto, index) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
                      >
                        {crypto.protocol}
                        {crypto.protocol_version && (
                          <span className="ml-1 text-blue-600 dark:text-blue-300">
                            {crypto.protocol_version}
                          </span>
                        )}
                      </span>
                    )) || (
                      <span className="text-sm text-gray-500 dark:text-gray-400">
                        No crypto detected
                      </span>
                    )}
                    {(asset.crypto_implementations?.length || 0) > 3 && (
                      <span className="text-xs text-gray-500 dark:text-gray-400">
                        +{(asset.crypto_implementations?.length || 0) - 3} more
                      </span>
                    )}
                  </div>
                </td>

                {/* Actions */}
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      onAssetClick?.(asset);
                    }}
                    className="text-primary-600 hover:text-primary-900 dark:text-primary-400 dark:hover:text-primary-300"
                  >
                    <EyeIcon className="h-5 w-5" />
                    <span className="sr-only">View details</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
