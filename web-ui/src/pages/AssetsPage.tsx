import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Asset, AssetFilters } from '../types/inventory';
import { inventoryApi } from '../services/inventoryApi';
import { AssetFiltersComponent } from '../components/inventory/AssetFilters';
import { AssetTable } from '../components/inventory/AssetTable';
import { AssetDetailsModal } from '../components/inventory/AssetDetailsModal';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import { Button } from '../components/common/Button';
import toast from 'react-hot-toast';

export const AssetsPage: React.FC = () => {
  const [filters, setFilters] = useState<AssetFilters>({
    page: 1,
    page_size: 20,
    sort_by: 'hostname',
    sort_order: 'asc',
  });
  
  const [selectedAsset, setSelectedAsset] = useState<Asset | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Fetch assets
  const {
    data: assetsResponse,
    isLoading,
    error,
    refetch
  } = useQuery({
    queryKey: ['assets', filters],
    queryFn: () => inventoryApi.getAssets(filters),
    retry: 2,
  });

  // Handle errors
  useEffect(() => {
    if (error) {
      toast.error('Failed to load assets. Please try again.');
      console.error('Assets query error:', error);
    }
  }, [error]);

  const handleFiltersChange = (newFilters: AssetFilters) => {
    setFilters(newFilters);
  };

  const handleClearFilters = () => {
    setFilters({
      page: 1,
      page_size: 20,
      sort_by: 'hostname',
      sort_order: 'asc',
    });
  };

  const handleAssetClick = async (asset: Asset) => {
    try {
      // If the asset doesn't have crypto implementations, fetch them
      if (!asset.crypto_implementations) {
        const detailedAsset = await inventoryApi.getAsset(asset.id);
        setSelectedAsset(detailedAsset);
      } else {
        setSelectedAsset(asset);
      }
      setIsModalOpen(true);
    } catch (error) {
      toast.error('Failed to load asset details.');
      console.error('Failed to fetch asset details:', error);
    }
  };

  const handlePageChange = (newPage: number) => {
    setFilters({ ...filters, page: newPage });
  };

  const assets = assetsResponse?.assets || [];
  const pagination = assetsResponse?.pagination;

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Crypto Asset Inventory
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Discover, analyze, and manage cryptographic implementations across your infrastructure
          </p>
        </div>

        {/* Filters */}
        <AssetFiltersComponent
          filters={filters}
          onFiltersChange={handleFiltersChange}
          onClearFilters={handleClearFilters}
        />

        {/* Results Summary */}
        {pagination && (
          <div className="mb-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
              <div className="text-sm text-gray-700 dark:text-gray-300">
                Showing{' '}
                <span className="font-medium">
                  {((pagination.page - 1) * pagination.page_size) + 1}
                </span>{' '}
                to{' '}
                <span className="font-medium">
                  {Math.min(pagination.page * pagination.page_size, pagination.total)}
                </span>{' '}
                of{' '}
                <span className="font-medium">{pagination.total}</span>{' '}
                assets
              </div>
              
              {/* Refresh Button */}
              <div className="mt-3 sm:mt-0">
                <Button
                  variant="secondary"
                  onClick={() => refetch()}
                  disabled={isLoading}
                  size="sm"
                >
                  {isLoading ? 'Refreshing...' : 'Refresh'}
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Assets Table */}
        <div className="mb-8">
          <AssetTable
            assets={assets}
            loading={isLoading}
            onAssetClick={handleAssetClick}
          />
        </div>

        {/* Pagination */}
        {pagination && pagination.total_pages > 1 && (
          <div className="flex items-center justify-between border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-4 py-3 sm:px-6 rounded-lg">
            <div className="flex flex-1 justify-between sm:hidden">
              <Button
                variant="secondary"
                onClick={() => handlePageChange(pagination.page - 1)}
                disabled={!pagination.has_prev}
              >
                Previous
              </Button>
              <Button
                variant="secondary"
                onClick={() => handlePageChange(pagination.page + 1)}
                disabled={!pagination.has_next}
              >
                Next
              </Button>
            </div>
            
            <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
              <div>
                <p className="text-sm text-gray-700 dark:text-gray-300">
                  Page{' '}
                  <span className="font-medium">{pagination.page}</span>{' '}
                  of{' '}
                  <span className="font-medium">{pagination.total_pages}</span>
                </p>
              </div>
              <div>
                <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                  <button
                    onClick={() => handlePageChange(pagination.page - 1)}
                    disabled={!pagination.has_prev}
                    className="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed dark:ring-gray-600 dark:hover:bg-gray-700"
                  >
                    <span className="sr-only">Previous</span>
                    <ChevronLeftIcon className="h-5 w-5" aria-hidden="true" />
                  </button>
                  
                  {/* Page numbers */}
                  {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                    let pageNum;
                    if (pagination.total_pages <= 5) {
                      pageNum = i + 1;
                    } else if (pagination.page <= 3) {
                      pageNum = i + 1;
                    } else if (pagination.page >= pagination.total_pages - 2) {
                      pageNum = pagination.total_pages - 4 + i;
                    } else {
                      pageNum = pagination.page - 2 + i;
                    }
                    
                    return (
                      <button
                        key={pageNum}
                        onClick={() => handlePageChange(pageNum)}
                        className={`relative inline-flex items-center px-4 py-2 text-sm font-semibold ${
                          pageNum === pagination.page
                            ? 'z-10 bg-primary-600 text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'
                            : 'text-gray-900 dark:text-gray-300 ring-1 ring-inset ring-gray-300 dark:ring-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 focus:z-20 focus:outline-offset-0'
                        }`}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                  
                  <button
                    onClick={() => handlePageChange(pagination.page + 1)}
                    disabled={!pagination.has_next}
                    className="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed dark:ring-gray-600 dark:hover:bg-gray-700"
                  >
                    <span className="sr-only">Next</span>
                    <ChevronRightIcon className="h-5 w-5" aria-hidden="true" />
                  </button>
                </nav>
              </div>
            </div>
          </div>
        )}

        {/* Asset Details Modal */}
        <AssetDetailsModal
          asset={selectedAsset}
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false);
            setSelectedAsset(null);
          }}
        />
      </div>
    </div>
  );
};
