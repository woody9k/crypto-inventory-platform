import React, { useState, useEffect } from 'react';
import { 
  DocumentTextIcon,
  ChartBarIcon,
  ShieldCheckIcon,
  NetworkIcon,
  ExclamationTriangleIcon,
  ArrowDownTrayIcon,
  EyeIcon,
  TrashIcon,
  PlusIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon
} from '@heroicons/react/24/outline';

/**
 * Report interface representing a generated report with metadata and status.
 * This matches the structure returned by the report generator API.
 */
interface Report {
  id: string;                    // Unique report identifier
  title: string;                 // Human-readable report title
  type: string;                  // Report type (crypto_summary, compliance_status, etc.)
  status: 'generating' | 'completed' | 'failed'; // Current report status
  created_at: string;            // ISO timestamp of report creation
  completed_at?: string;         // ISO timestamp of report completion (optional)
  download_url?: string;         // URL for downloading the report (optional)
}

/**
 * ReportTemplate interface representing a predefined report template.
 * Templates define the structure and parameters for different report types.
 */
interface ReportTemplate {
  id: string;          // Unique template identifier
  name: string;        // Human-readable template name
  description: string; // Template description
  type: string;        // Template type (summary, compliance, etc.)
  category: string;    // Template category (crypto, compliance, network, security)
}

interface ReportsPageProps {}

/**
 * ReportsPage component provides a comprehensive interface for managing reports.
 * Features include:
 * - Viewing all generated reports with status tracking
 * - Generating new reports from templates
 * - Downloading completed reports
 * - Deleting reports
 * - Real-time status updates
 */
const ReportsPage: React.FC<ReportsPageProps> = () => {
  const [reports, setReports] = useState<Report[]>([]);
  const [templates, setTemplates] = useState<ReportTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [showGenerateModal, setShowGenerateModal] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<ReportTemplate | null>(null);

  // Load mock data for demonstration purposes
  // In production, this would make API calls to fetch real data
  useEffect(() => {
    const mockReports: Report[] = [
      {
        id: 'report-001',
        title: 'Crypto Summary Report',
        type: 'crypto_summary',
        status: 'completed',
        created_at: '2024-12-15T09:00:00Z',
        completed_at: '2024-12-15T09:02:30Z',
        download_url: '/api/v1/reports/report-001/download'
      },
      {
        id: 'report-002',
        title: 'Compliance Status Report',
        type: 'compliance_status',
        status: 'completed',
        created_at: '2024-12-15T08:30:00Z',
        completed_at: '2024-12-15T08:32:15Z',
        download_url: '/api/v1/reports/report-002/download'
      },
      {
        id: 'report-003',
        title: 'Network Topology Report',
        type: 'network_topology',
        status: 'generating',
        created_at: '2024-12-15T10:15:00Z'
      },
      {
        id: 'report-004',
        title: 'Risk Assessment Report',
        type: 'risk_assessment',
        status: 'failed',
        created_at: '2024-12-15T07:45:00Z'
      }
    ];

    const mockTemplates: ReportTemplate[] = [
      {
        id: 'crypto_summary',
        name: 'Crypto Summary Report',
        description: 'Overview of all cryptographic implementations across the network',
        type: 'summary',
        category: 'crypto'
      },
      {
        id: 'compliance_status',
        name: 'Compliance Status Report',
        description: 'Current compliance status against various frameworks',
        type: 'compliance',
        category: 'compliance'
      },
      {
        id: 'network_topology',
        name: 'Network Topology Report',
        description: 'Network topology and sensor coverage map',
        type: 'topology',
        category: 'network'
      },
      {
        id: 'risk_assessment',
        name: 'Risk Assessment Report',
        description: 'Security risk assessment and recommendations',
        type: 'risk',
        category: 'security'
      },
      {
        id: 'certificate_audit',
        name: 'Certificate Audit Report',
        description: 'SSL/TLS certificate inventory and expiration analysis',
        type: 'audit',
        category: 'crypto'
      }
    ];

    setTimeout(() => {
      setReports(mockReports);
      setTemplates(mockTemplates);
      setLoading(false);
    }, 1000);
  }, []);

  /**
   * Returns the appropriate icon component for a given report status.
   * @param status - The report status (generating, completed, failed)
   * @returns JSX element representing the status icon
   */
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircleIcon className="h-5 w-5 text-green-500" />;
      case 'generating':
        return <ClockIcon className="h-5 w-5 text-yellow-500" />;
      case 'failed':
        return <XCircleIcon className="h-5 w-5 text-red-500" />;
      default:
        return <ClockIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  /**
   * Returns the appropriate CSS classes for a given report status.
   * @param status - The report status (generating, completed, failed)
   * @returns String of CSS classes for styling the status badge
   */
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'generating':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  /**
   * Returns the appropriate icon component for a given report type.
   * @param type - The report type (crypto_summary, compliance_status, etc.)
   * @returns JSX element representing the report type icon
   */
  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'crypto_summary':
        return <ChartBarIcon className="h-6 w-6 text-blue-500" />;
      case 'compliance_status':
        return <ShieldCheckIcon className="h-6 w-6 text-green-500" />;
      case 'network_topology':
        return <NetworkIcon className="h-6 w-6 text-purple-500" />;
      case 'risk_assessment':
        return <ExclamationTriangleIcon className="h-6 w-6 text-red-500" />;
      default:
        return <DocumentTextIcon className="h-6 w-6 text-gray-500" />;
    }
  };

  /**
   * Formats a date string into a human-readable format.
   * @param dateString - ISO date string
   * @returns Formatted date string
   */
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  /**
   * Handles the generation of a new report from a template.
   * @param template - The selected report template
   */
  const handleGenerateReport = (template: ReportTemplate) => {
    setSelectedTemplate(template);
    setShowGenerateModal(true);
  };

  /**
   * Handles the download of a completed report.
   * @param report - The report to download
   */
  const handleDownloadReport = (report: Report) => {
    if (report.download_url) {
      // In a real implementation, this would trigger the download
      console.log('Downloading report:', report.download_url);
    }
  };

  /**
   * Handles the deletion of a report from the system.
   * @param reportId - The ID of the report to delete
   */
  const handleDeleteReport = (reportId: string) => {
    setReports(reports.filter(r => r.id !== reportId));
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Reports
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Generate and manage compliance and security reports
          </p>
        </div>
        <button
          onClick={() => setShowGenerateModal(true)}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <PlusIcon className="h-4 w-4 mr-2" />
          Generate Report
        </button>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <DocumentTextIcon className="h-6 w-6 text-blue-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Total Reports
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {reports.length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <CheckCircleIcon className="h-6 w-6 text-green-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Completed
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {reports.filter(r => r.status === 'completed').length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <ClockIcon className="h-6 w-6 text-yellow-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Generating
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {reports.filter(r => r.status === 'generating').length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <XCircleIcon className="h-6 w-6 text-red-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Failed
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {reports.filter(r => r.status === 'failed').length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Reports List */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
            Recent Reports
          </h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500 dark:text-gray-400">
            View and manage your generated reports
          </p>
        </div>
        <ul className="divide-y divide-gray-200 dark:divide-gray-700">
          {reports.map((report) => (
            <li key={report.id}>
              <div className="px-4 py-4 sm:px-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    {getTypeIcon(report.type)}
                    <div className="ml-4">
                      <div className="flex items-center">
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {report.title}
                        </p>
                        <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(report.status)}`}>
                          {report.status}
                        </span>
                      </div>
                      <div className="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                        <p className="mr-4">Created: {formatDate(report.created_at)}</p>
                        {report.completed_at && (
                          <p>Completed: {formatDate(report.completed_at)}</p>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    {report.status === 'completed' && (
                      <button
                        onClick={() => handleDownloadReport(report)}
                        className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                      >
                        <ArrowDownTrayIcon className="h-4 w-4 mr-1" />
                        Download
                      </button>
                    )}
                    <button
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <EyeIcon className="h-4 w-4 mr-1" />
                      View
                    </button>
                    <button
                      onClick={() => handleDeleteReport(report.id)}
                      className="inline-flex items-center px-3 py-1 border border-red-300 dark:border-red-600 text-sm font-medium rounded-md text-red-700 dark:text-red-300 bg-white dark:bg-gray-700 hover:bg-red-50 dark:hover:bg-red-600"
                    >
                      <TrashIcon className="h-4 w-4 mr-1" />
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Generate Report Modal */}
      {showGenerateModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white dark:bg-gray-800">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                Generate New Report
              </h3>
              <div className="space-y-4">
                {templates.map((template) => (
                  <div
                    key={template.id}
                    className="border border-gray-200 dark:border-gray-600 rounded-lg p-4 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700"
                    onClick={() => handleGenerateReport(template)}
                  >
                    <div className="flex items-center">
                      {getTypeIcon(template.id)}
                      <div className="ml-3">
                        <h4 className="text-sm font-medium text-gray-900 dark:text-white">
                          {template.name}
                        </h4>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {template.description}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
              <div className="flex justify-end space-x-2 mt-6">
                <button
                  onClick={() => setShowGenerateModal(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ReportsPage;
