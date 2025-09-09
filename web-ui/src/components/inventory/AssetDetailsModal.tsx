import React from 'react';
import { XMarkIcon, ShieldCheckIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline';
import { Asset, CryptoImplementation } from '../../types/inventory';
import { RiskBadge } from './RiskBadge';

interface AssetDetailsModalProps {
  asset: Asset | null;
  isOpen: boolean;
  onClose: () => void;
}

export const AssetDetailsModal: React.FC<AssetDetailsModalProps> = ({
  asset,
  isOpen,
  onClose,
}) => {
  if (!isOpen || !asset) return null;

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const CryptoImplementationCard: React.FC<{ crypto: CryptoImplementation }> = ({ crypto }) => (
    <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center space-x-2">
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
            {crypto.protocol}
            {crypto.protocol_version && ` ${crypto.protocol_version}`}
          </span>
          <RiskBadge riskLevel={crypto.risk_level} riskScore={crypto.risk_score} size="sm" />
        </div>
        <span className="text-xs text-gray-500 dark:text-gray-400 capitalize">
          {crypto.discovery_method}
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
        {crypto.cipher_suite && (
          <div>
            <span className="text-gray-500 dark:text-gray-400">Cipher Suite:</span>
            <div className="font-mono text-xs bg-white dark:bg-gray-900 rounded px-2 py-1 mt-1">
              {crypto.cipher_suite}
            </div>
          </div>
        )}
        
        {crypto.key_size && (
          <div>
            <span className="text-gray-500 dark:text-gray-400">Key Size:</span>
            <div className="font-medium">{crypto.key_size} bits</div>
          </div>
        )}

        {crypto.hash_algorithm && (
          <div>
            <span className="text-gray-500 dark:text-gray-400">Hash Algorithm:</span>
            <div className="font-medium">{crypto.hash_algorithm}</div>
          </div>
        )}

        {crypto.signature_algorithm && (
          <div>
            <span className="text-gray-500 dark:text-gray-400">Signature:</span>
            <div className="font-medium">{crypto.signature_algorithm}</div>
          </div>
        )}
      </div>

      {crypto.risk_factors && crypto.risk_factors.length > 0 && (
        <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-600">
          <div className="flex items-center text-orange-600 dark:text-orange-400 mb-2">
            <ExclamationTriangleIcon className="h-4 w-4 mr-1" />
            <span className="text-sm font-medium">Risk Factors</span>
          </div>
          <div className="flex flex-wrap gap-1">
            {crypto.risk_factors.map((factor, index) => (
              <span
                key={index}
                className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200"
              >
                {factor}
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-600 text-xs text-gray-500 dark:text-gray-400">
        <div className="flex justify-between">
          <span>Confidence: {Math.round(crypto.confidence_score * 100)}%</span>
          <span>Last Verified: {formatDate(crypto.last_verified_at)}</span>
        </div>
      </div>
    </div>
  );

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        {/* Background overlay */}
        <div
          className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          onClick={onClose}
        />

        {/* Modal panel */}
        <div className="inline-block align-bottom bg-white dark:bg-gray-900 rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full sm:p-6">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-2xl font-semibold text-gray-900 dark:text-white">
                {asset.hostname || 'Unknown Asset'}
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400 capitalize">
                {asset.asset_type}
                {asset.operating_system && ` • ${asset.operating_system}`}
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <RiskBadge riskLevel={asset.risk_level} riskScore={asset.risk_score} />
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <XMarkIcon className="h-6 w-6" />
              </button>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Asset Information */}
            <div className="space-y-6">
              <div>
                <h4 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                  Asset Information
                </h4>
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 space-y-3">
                  <div className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">IP Address:</span>
                      <div className="font-medium">{asset.ip_address || 'N/A'}</div>
                    </div>
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">Port:</span>
                      <div className="font-medium">{asset.port || 'N/A'}</div>
                    </div>
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">Environment:</span>
                      <div className="font-medium capitalize">{asset.environment || 'N/A'}</div>
                    </div>
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">Business Unit:</span>
                      <div className="font-medium">{asset.business_unit || 'N/A'}</div>
                    </div>
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">Owner:</span>
                      <div className="font-medium">{asset.owner_email || 'N/A'}</div>
                    </div>
                    <div>
                      <span className="text-gray-500 dark:text-gray-400">Asset Type:</span>
                      <div className="font-medium capitalize">{asset.asset_type}</div>
                    </div>
                  </div>
                  
                  {asset.description && (
                    <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                      <span className="text-gray-500 dark:text-gray-400">Description:</span>
                      <div className="font-medium mt-1">{asset.description}</div>
                    </div>
                  )}
                </div>
              </div>

              {/* Timestamps */}
              <div>
                <h4 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                  Timeline
                </h4>
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 space-y-3 text-sm">
                  <div>
                    <span className="text-gray-500 dark:text-gray-400">First Discovered:</span>
                    <div className="font-medium">{formatDate(asset.first_discovered_at)}</div>
                  </div>
                  <div>
                    <span className="text-gray-500 dark:text-gray-400">Last Seen:</span>
                    <div className="font-medium">{formatDate(asset.last_seen_at)}</div>
                  </div>
                  <div>
                    <span className="text-gray-500 dark:text-gray-400">Created:</span>
                    <div className="font-medium">{formatDate(asset.created_at)}</div>
                  </div>
                  <div>
                    <span className="text-gray-500 dark:text-gray-400">Updated:</span>
                    <div className="font-medium">{formatDate(asset.updated_at)}</div>
                  </div>
                </div>
              </div>

              {/* Tags */}
              {asset.tags && Object.keys(asset.tags).length > 0 && (
                <div>
                  <h4 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                    Tags
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {Object.entries(asset.tags).map(([key, value]) => (
                      <span
                        key={key}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
                      >
                        {key}: {String(value)}
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Crypto Implementations */}
            <div>
              <div className="flex items-center mb-4">
                <ShieldCheckIcon className="h-5 w-5 text-blue-600 dark:text-blue-400 mr-2" />
                <h4 className="text-lg font-medium text-gray-900 dark:text-white">
                  Cryptographic Implementations
                </h4>
                <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">
                  ({asset.crypto_implementations?.length || 0})
                </span>
              </div>

              <div className="space-y-4 max-h-96 overflow-y-auto">
                {asset.crypto_implementations && asset.crypto_implementations.length > 0 ? (
                  asset.crypto_implementations.map((crypto) => (
                    <CryptoImplementationCard key={crypto.id} crypto={crypto} />
                  ))
                ) : (
                  <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    <ShieldCheckIcon className="h-12 w-12 mx-auto mb-3 opacity-50" />
                    <p>No cryptographic implementations detected</p>
                    <p className="text-sm">This asset may not be using crypto protocols or hasn't been scanned yet.</p>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Footer */}
          <div className="mt-8 pt-6 border-t border-gray-200 dark:border-gray-700">
            <div className="flex justify-end">
              <button
                onClick={onClose}
                className="px-4 py-2 bg-gray-300 hover:bg-gray-400 dark:bg-gray-600 dark:hover:bg-gray-500 text-gray-700 dark:text-gray-200 rounded-lg transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
