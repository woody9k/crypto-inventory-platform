import React, { useState, useEffect } from 'react';
import { 
  PlusIcon, 
  TrashIcon, 
  PlayIcon, 
  StopIcon, 
  CogIcon,
  EyeIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  ClockIcon,
  ClipboardDocumentIcon,
  XCircleIcon
} from '@heroicons/react/24/outline';

interface Sensor {
  id: string;
  name: string;
  status: 'active' | 'inactive' | 'error' | 'unknown' | 'pending';
  platform: string;
  version: string;
  profile: string;
  lastSeen: string;
  ipAddress: string;
  networkInterfaces: string[];
  uptime: number;
  memoryUsage: number;
  cpuUsage: number;
  packetsCaptured: number;
  discoveriesMade: number;
  errors: number;
  registrationKey?: string;
  expiresAt?: string;
}

interface SensorManagementPageProps {}

const SensorManagementPage: React.FC<SensorManagementPageProps> = () => {
  const [sensors, setSensors] = useState<Sensor[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddSensor, setShowAddSensor] = useState(false);
  const [, setSelectedSensor] = useState<Sensor | null>(null);

  // Delete sensor function
  const deleteSensor = (id: string) => {
    setSensors(sensors.filter(s => s.id !== id));
  };

  // Mock data for demonstration
  useEffect(() => {
    const mockSensors: Sensor[] = [
      {
        id: 'sensor-dc01-eth0-20241215',
        name: 'sensor-dc01',
        status: 'active',
        platform: 'linux',
        version: '1.0.0',
        profile: 'datacenter_host',
        lastSeen: '2024-12-15T10:30:00Z',
        ipAddress: '192.168.1.100',
        networkInterfaces: ['eth0', 'eth1'],
        uptime: 3600,
        memoryUsage: 52428800,
        cpuUsage: 15.5,
        packetsCaptured: 15000,
        discoveriesMade: 45,
        errors: 0
      },
      {
        id: 'sensor-cloud01-ens3-20241215',
        name: 'sensor-cloud01',
        status: 'active',
        platform: 'linux',
        version: '1.0.0',
        profile: 'cloud_instance',
        lastSeen: '2024-12-15T10:29:45Z',
        ipAddress: '10.0.1.50',
        networkInterfaces: ['ens3'],
        uptime: 7200,
        memoryUsage: 31457280,
        cpuUsage: 8.2,
        packetsCaptured: 8500,
        discoveriesMade: 23,
        errors: 1
      },
      {
        id: 'sensor-user01-wlan0-20241215',
        name: 'sensor-user01',
        status: 'error',
        platform: 'linux',
        version: '1.0.0',
        profile: 'end_user_machine',
        lastSeen: '2024-12-15T09:15:30Z',
        ipAddress: '192.168.0.25',
        networkInterfaces: ['wlan0'],
        uptime: 0,
        memoryUsage: 0,
        cpuUsage: 0,
        packetsCaptured: 0,
        discoveriesMade: 0,
        errors: 5
      },
      {
        id: 'pending-sensor-dc02-20241215',
        name: 'sensor-dc02',
        status: 'pending',
        platform: 'linux',
        version: '1.0.0',
        profile: 'datacenter_host',
        lastSeen: '2024-12-15T10:00:00Z',
        ipAddress: '192.168.1.101',
        networkInterfaces: ['eth0', 'eth1'],
        uptime: 0,
        memoryUsage: 0,
        cpuUsage: 0,
        packetsCaptured: 0,
        discoveriesMade: 0,
        errors: 0,
        registrationKey: 'REG-tenant-123-20241215-D0E6F2',
        expiresAt: '2024-12-15T11:00:00Z'
      }
    ];

    setTimeout(() => {
      setSensors(mockSensors);
      setLoading(false);
    }, 1000);
  }, []);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'active':
        return <CheckCircleIcon className="h-5 w-5 text-green-500" />;
      case 'error':
        return <XCircleIcon className="h-5 w-5 text-red-500" />;
      case 'inactive':
        return <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500" />;
      case 'pending':
        return <ClockIcon className="h-5 w-5 text-blue-500" />;
      default:
        return <ExclamationTriangleIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'error':
        return 'bg-red-100 text-red-800';
      case 'inactive':
        return 'bg-yellow-100 text-yellow-800';
      case 'pending':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatUptime = (seconds: number) => {
    if (seconds === 0) return 'Offline';
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  const formatTimeRemaining = (expiresAt: string) => {
    const now = new Date().getTime();
    const expiry = new Date(expiresAt).getTime();
    const diff = expiry - now;

    if (diff <= 0) return 'Expired';

    const minutes = Math.floor(diff / 60000);
    const seconds = Math.floor((diff % 60000) / 1000);
    return `${minutes}m ${seconds}s`;
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
            Sensor Management
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Manage network sensors and monitor their status
          </p>
        </div>
        <button
          onClick={() => setShowAddSensor(true)}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <PlusIcon className="h-4 w-4 mr-2" />
          Add Sensor
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <CheckCircleIcon className="h-6 w-6 text-green-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Active Sensors
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {sensors.filter(s => s.status === 'active').length}
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
                <ExclamationTriangleIcon className="h-6 w-6 text-yellow-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Inactive Sensors
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {sensors.filter(s => s.status === 'inactive').length}
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
                    Error Sensors
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {sensors.filter(s => s.status === 'error').length}
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
                <EyeIcon className="h-6 w-6 text-blue-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Total Discoveries
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {sensors.reduce((sum, s) => sum + s.discoveriesMade, 0)}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Sensors Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
            Network Sensors
          </h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500 dark:text-gray-400">
            Manage and monitor your network sensors
          </p>
        </div>
        <ul className="divide-y divide-gray-200 dark:divide-gray-700">
          {sensors.map((sensor) => (
            <li key={sensor.id}>
              <div className="px-4 py-4 sm:px-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    {getStatusIcon(sensor.status)}
                    <div className="ml-4">
                      <div className="flex items-center">
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {sensor.name}
                        </p>
                        <span className={`ml-2 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(sensor.status)}`}>
                          {sensor.status}
                        </span>
                      </div>
                      <div className="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                        <p className="mr-4">Platform: {sensor.platform}</p>
                        <p className="mr-4">Version: {sensor.version}</p>
                        <p className="mr-4">Profile: {sensor.profile}</p>
                        <p className="mr-4">IP: {sensor.ipAddress}</p>
                        {sensor.status === 'pending' && sensor.expiresAt && (
                          <p className="text-blue-600 dark:text-blue-400">
                            Expires in: {formatTimeRemaining(sensor.expiresAt)}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => setSelectedSensor(sensor)}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <EyeIcon className="h-4 w-4 mr-1" />
                      View
                    </button>
                    <button
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <CogIcon className="h-4 w-4 mr-1" />
                      Config
                    </button>
                    {sensor.status === 'active' ? (
                      <button className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                        <StopIcon className="h-4 w-4 mr-1" />
                        Stop
                      </button>
                    ) : (
                      <button className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                        <PlayIcon className="h-4 w-4 mr-1" />
                        Start
                      </button>
                    )}
                    <button 
                      onClick={() => deleteSensor(sensor.id)}
                      className="inline-flex items-center px-3 py-1 border border-red-300 dark:border-red-600 text-sm font-medium rounded-md text-red-700 dark:text-red-300 bg-white dark:bg-gray-700 hover:bg-red-50 dark:hover:bg-red-600"
                    >
                      <TrashIcon className="h-4 w-4 mr-1" />
                      Remove
                    </button>
                  </div>
                </div>
                
                {/* Sensor Details */}
                <div className="mt-4 grid grid-cols-1 md:grid-cols-4 gap-4">
                  <div>
                    <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                      Network Interfaces
                    </p>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {sensor.networkInterfaces.join(', ')}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                      Uptime
                    </p>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {formatUptime(sensor.uptime)}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                      Memory Usage
                    </p>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {formatBytes(sensor.memoryUsage)}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                      CPU Usage
                    </p>
                    <p className="text-sm text-gray-900 dark:text-white">
                      {sensor.cpuUsage}%
                    </p>
                  </div>
                </div>

                {/* Pending Sensor Information */}
                {sensor.status === 'pending' && sensor.registrationKey && (
                  <div className="mt-4 bg-blue-50 dark:bg-blue-900/20 rounded-md p-4">
                    <h4 className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-2">
                      Installation Instructions
                    </h4>
                    <div className="space-y-2">
                      <div>
                        <p className="text-xs font-medium text-blue-700 dark:text-blue-300 uppercase tracking-wide">
                          Registration Key
                        </p>
                        <div className="flex items-center space-x-2">
                          <code className="text-sm text-blue-900 dark:text-blue-100 bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded">
                            {sensor.registrationKey}
                          </code>
                          <button
                            onClick={() => navigator.clipboard.writeText(sensor.registrationKey!)}
                            className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200"
                          >
                            <ClipboardDocumentIcon className="h-4 w-4" />
                          </button>
                        </div>
                      </div>
                      <div>
                        <p className="text-xs font-medium text-blue-700 dark:text-blue-300 uppercase tracking-wide">
                          Installation Command
                        </p>
                        <div className="flex items-center space-x-2">
                          <code className="text-sm text-blue-900 dark:text-blue-100 bg-blue-100 dark:bg-blue-800 px-2 py-1 rounded flex-1">
                            curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key {sensor.registrationKey} --ip {sensor.ipAddress} --name "{sensor.name}" --profile {sensor.profile}
                          </code>
                          <button
                            onClick={() => navigator.clipboard.writeText(`curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key ${sensor.registrationKey} --ip ${sensor.ipAddress} --name "${sensor.name}" --profile ${sensor.profile}`)}
                            className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200"
                          >
                            <ClipboardDocumentIcon className="h-4 w-4" />
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Performance Metrics */}
                {sensor.status !== 'pending' && (
                  <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                        Packets Captured
                      </p>
                      <p className="text-sm text-gray-900 dark:text-white">
                        {sensor.packetsCaptured.toLocaleString()}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                        Discoveries Made
                      </p>
                      <p className="text-sm text-gray-900 dark:text-white">
                        {sensor.discoveriesMade}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                        Errors
                      </p>
                      <p className="text-sm text-gray-900 dark:text-white">
                        {sensor.errors}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Add Sensor Modal */}
      {showAddSensor && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white dark:bg-gray-800">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                Add New Sensor
              </h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Registration Key
                  </label>
                  <input
                    type="text"
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="REG-550e8400-20241215-A7B3C9"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Sensor Name
                  </label>
                  <input
                    type="text"
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="sensor-dc01"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Network Interfaces
                  </label>
                  <input
                    type="text"
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="eth0,eth1"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Profile
                  </label>
                  <select className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white">
                    <option value="datacenter_host">Datacenter Host</option>
                    <option value="cloud_instance">Cloud Instance</option>
                    <option value="end_user_machine">End User Machine</option>
                    <option value="air_gapped">Air Gapped</option>
                  </select>
                </div>
              </div>
              <div className="flex justify-end space-x-2 mt-6">
                <button
                  onClick={() => setShowAddSensor(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
                >
                  Cancel
                </button>
                <button
                  onClick={() => setShowAddSensor(false)}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  Add Sensor
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SensorManagementPage;
