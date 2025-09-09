import React, { useState, useEffect } from 'react';
import { 
  PlusIcon, 
  ClockIcon, 
  CheckCircleIcon, 
  ExclamationTriangleIcon,
  ClipboardDocumentIcon,
  TrashIcon,
  CogIcon,
} from '@heroicons/react/24/outline';

interface PendingSensor {
  id: string;
  name: string;
  ipAddress: string;
  tags: string[];
  profile: string;
  networkInterfaces: string[];
  registrationKey: string;
  createdAt: string;
  expiresAt: string;
  status: 'pending' | 'expired' | 'used';
  installationCommand: string;
}

interface AdminSettings {
  keyExpirationMinutes: number;
  maxPendingSensors: number;
  requireIpValidation: boolean;
}

const SensorRegistrationPage: React.FC = () => {
  const [pendingSensors, setPendingSensors] = useState<PendingSensor[]>([]);
  const [adminSettings, setAdminSettings] = useState<AdminSettings>({
    keyExpirationMinutes: 60,
    maxPendingSensors: 50,
    requireIpValidation: true
  });
  const [showAddSensor, setShowAddSensor] = useState(false);
  const [showAdminSettings, setShowAdminSettings] = useState(false);
  const [loading, setLoading] = useState(true);

  // Form state for new sensor
  const [formData, setFormData] = useState({
    name: '',
    ipAddress: '',
    tags: [] as string[],
    profile: 'datacenter_host',
    networkInterfaces: [] as string[],
    description: ''
  });

  // Available tags and profiles
  const availableTags = [
    'datacenter', 'office', 'cloud', 'critical', 'test', 'production',
    'staging', 'development', 'dmz', 'internal', 'external'
  ];

  const availableProfiles = [
    { value: 'datacenter_host', label: 'Datacenter Host', description: 'Full features, multiple interfaces' },
    { value: 'cloud_instance', label: 'Cloud Instance', description: 'Cloud-optimized, periodic reporting' },
    { value: 'end_user_machine', label: 'End User Machine', description: 'Minimal footprint, single interface' },
    { value: 'air_gapped', label: 'Air Gapped', description: 'Offline mode, export files' }
  ];

  // Mock data for demonstration
  useEffect(() => {
    const mockPendingSensors: PendingSensor[] = [
      {
        id: 'pending-001',
        name: 'sensor-dc01',
        ipAddress: '192.168.1.100',
        tags: ['datacenter', 'critical'],
        profile: 'datacenter_host',
        networkInterfaces: ['eth0', 'eth1'],
        registrationKey: 'REG-tenant-123-20241215-A7B3C9',
        createdAt: '2024-12-15T10:00:00Z',
        expiresAt: '2024-12-15T11:00:00Z',
        status: 'pending',
        installationCommand: 'curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key REG-tenant-123-20241215-A7B3C9 --ip 192.168.1.100 --name "Datacenter Sensor 1" --profile datacenter_host'
      },
      {
        id: 'pending-002',
        name: 'sensor-cloud01',
        ipAddress: '10.0.1.50',
        tags: ['cloud', 'production'],
        profile: 'cloud_instance',
        networkInterfaces: ['ens3'],
        registrationKey: 'REG-tenant-123-20241215-B8C4D0',
        createdAt: '2024-12-15T09:30:00Z',
        expiresAt: '2024-12-15T10:30:00Z',
        status: 'pending',
        installationCommand: 'curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key REG-tenant-123-20241215-B8C4D0 --ip 10.0.1.50 --name "Cloud Sensor 1" --profile cloud_vm'
      },
      {
        id: 'pending-003',
        name: 'sensor-user01',
        ipAddress: '192.168.0.25',
        tags: ['office', 'test'],
        profile: 'end_user_machine',
        networkInterfaces: ['wlan0'],
        registrationKey: 'REG-tenant-123-20241215-C9D5E1',
        createdAt: '2024-12-15T08:00:00Z',
        expiresAt: '2024-12-15T09:00:00Z',
        status: 'expired',
        installationCommand: 'curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key REG-tenant-123-20241215-C9D5E1 --ip 192.168.0.25 --name "Edge Sensor 1" --profile edge_device'
      }
    ];

    setTimeout(() => {
      setPendingSensors(mockPendingSensors);
      setLoading(false);
    }, 1000);
  }, []);

  const formatTimeRemaining = (expiresAt: string) => {
    const now = new Date().getTime();
    const expiry = new Date(expiresAt).getTime();
    const diff = expiry - now;

    if (diff <= 0) return 'Expired';

    const minutes = Math.floor(diff / 60000);
    const seconds = Math.floor((diff % 60000) / 1000);
    return `${minutes}m ${seconds}s`;
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <ClockIcon className="h-5 w-5 text-yellow-500" />;
      case 'used':
        return <CheckCircleIcon className="h-5 w-5 text-green-500" />;
      case 'expired':
        return <ExclamationTriangleIcon className="h-5 w-5 text-red-500" />;
      default:
        return <ExclamationTriangleIcon className="h-5 w-5 text-gray-500" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'used':
        return 'bg-green-100 text-green-800';
      case 'expired':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  /**
   * Generates installation commands for a pending sensor
   * Creates three different installation methods:
   * 1. One-line curl command (recommended)
   * 2. Interactive mode for guided installation
   * 3. Manual download and installation
   * 
   * @param sensor - The pending sensor configuration
   * @returns Object containing all three installation command variants
   */
  const generateInstallationCommands = (sensor: PendingSensor) => {
    const interfaces = sensor.networkInterfaces.length > 0 ? sensor.networkInterfaces.join(',') : 'eth0';
    const controlPlaneURL = 'https://crypto-inventory.company.com';
    
    return {
      // Direct download and run - one-line installation
      curlCommand: `curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key ${sensor.registrationKey} --ip ${sensor.ipAddress} --name "${sensor.name}" --profile ${sensor.profile}`,
      
      // Interactive mode - guided installation with prompts
      interactiveCommand: `curl -sSL ${controlPlaneURL}/scripts/install-sensor.sh | sudo bash -s -- --interactive`,
      
      // Manual download - for air-gapped or restricted environments
      manualCommand: `# Download and run manually
curl -O ${controlPlaneURL}/scripts/install-sensor.sh
chmod +x install-sensor.sh
sudo ./install-sensor.sh \\
  --key ${sensor.registrationKey} \\
  --ip ${sensor.ipAddress} \\
  --name ${sensor.name} \\
  --profile ${sensor.profile} \\
  --interfaces "${interfaces}" \\
  --url ${controlPlaneURL}`
    };
  };

  const handleAddSensor = () => {
    // Validate form
    if (!formData.name || !formData.ipAddress) {
      alert('Please fill in required fields');
      return;
    }

    // Generate registration key
    const keySuffix = Math.random().toString(36).substring(2, 8).toUpperCase();
    const now = new Date();
    const key = `REG-tenant-123-${now.toISOString().slice(0, 10).replace(/-/g, '')}-${keySuffix}`;
    
    // Calculate expiration
    const expiresAt = new Date(now.getTime() + adminSettings.keyExpirationMinutes * 60000);

    // Create pending sensor
    const newSensor: PendingSensor = {
      id: `pending-${Date.now()}`,
      name: formData.name,
      ipAddress: formData.ipAddress,
      tags: formData.tags,
      profile: formData.profile,
      networkInterfaces: formData.networkInterfaces,
      registrationKey: key,
      createdAt: now.toISOString(),
      expiresAt: expiresAt.toISOString(),
      status: 'pending',
      installationCommand: `curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key ${key} --ip ${formData.ipAddress} --name "${formData.name}" --profile ${formData.profile}`
    };

    setPendingSensors([...pendingSensors, newSensor]);
    setShowAddSensor(false);
    setFormData({
      name: '',
      ipAddress: '',
      tags: [],
      profile: 'datacenter_host',
      networkInterfaces: [],
      description: ''
    });
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // You could add a toast notification here
  };

  const deletePendingSensor = (id: string) => {
    setPendingSensors(pendingSensors.filter(s => s.id !== id));
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
            Sensor Registration
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Create registration keys for new sensors
          </p>
        </div>
        <div className="flex space-x-3">
          <button
            onClick={() => setShowAdminSettings(true)}
            className="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
          >
            <CogIcon className="h-4 w-4 mr-2" />
            Admin Settings
          </button>
          <button
            onClick={() => setShowAddSensor(true)}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <PlusIcon className="h-4 w-4 mr-2" />
            Add Sensor
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <ClockIcon className="h-6 w-6 text-yellow-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Pending Sensors
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {pendingSensors.filter(s => s.status === 'pending').length}
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
                    Registered Today
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {pendingSensors.filter(s => s.status === 'used').length}
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
                <ExclamationTriangleIcon className="h-6 w-6 text-red-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Expired Keys
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {pendingSensors.filter(s => s.status === 'expired').length}
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
                <ClockIcon className="h-6 w-6 text-blue-500" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                    Key Expiration
                  </dt>
                  <dd className="text-lg font-medium text-gray-900 dark:text-white">
                    {adminSettings.keyExpirationMinutes}m
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Pending Sensors Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
            Pending Sensor Registrations
          </h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500 dark:text-gray-400">
            Registration keys waiting for sensor installation
          </p>
        </div>
        <ul className="divide-y divide-gray-200 dark:divide-gray-700">
          {pendingSensors.map((sensor) => (
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
                        {sensor.status === 'pending' && (
                          <span className="ml-2 text-sm text-gray-500 dark:text-gray-400">
                            Expires in: {formatTimeRemaining(sensor.expiresAt)}
                          </span>
                        )}
                      </div>
                      <div className="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                        <p className="mr-4">IP: {sensor.ipAddress}</p>
                        <p className="mr-4">Profile: {sensor.profile}</p>
                        <p className="mr-4">Tags: {sensor.tags.join(', ')}</p>
                        <p>Key: {sensor.registrationKey}</p>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => copyToClipboard(sensor.registrationKey)}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4 mr-1" />
                      Copy Key
                    </button>
                    <button
                      onClick={() => copyToClipboard(sensor.installationCommand)}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4 mr-1" />
                      Copy Command
                    </button>
                    <button
                      onClick={() => deletePendingSensor(sensor.id)}
                      className="inline-flex items-center px-3 py-1 border border-red-300 dark:border-red-600 text-sm font-medium rounded-md text-red-700 dark:text-red-300 bg-white dark:bg-gray-700 hover:bg-red-50 dark:hover:bg-red-600"
                    >
                      <TrashIcon className="h-4 w-4 mr-1" />
                      Delete
                    </button>
                  </div>
                </div>
                
                {/* Installation Commands */}
                <div className="mt-4">
                  <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
                    Installation Commands
                  </p>
                  
                  {(() => {
                    const commands = generateInstallationCommands(sensor);
                    return (
                      <div className="space-y-3">
                        {/* One-line curl command */}
                        <div>
                          <p className="text-xs font-medium text-blue-700 dark:text-blue-300 mb-1">
                            🚀 One-Line Installation (Recommended)
                          </p>
                          <div className="bg-blue-50 dark:bg-blue-900/20 rounded-md p-3 border border-blue-200 dark:border-blue-800">
                            <div className="flex items-start justify-between">
                              <code className="text-sm text-blue-900 dark:text-blue-100 flex-1">
                                {commands.curlCommand}
                              </code>
                              <button
                                onClick={() => copyToClipboard(commands.curlCommand)}
                                className="ml-2 text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-200"
                              >
                                <ClipboardDocumentIcon className="h-4 w-4" />
                              </button>
                            </div>
                          </div>
                        </div>

                        {/* Interactive mode */}
                        <div>
                          <p className="text-xs font-medium text-green-700 dark:text-green-300 mb-1">
                            🎛️ Interactive Mode
                          </p>
                          <div className="bg-green-50 dark:bg-green-900/20 rounded-md p-3 border border-green-200 dark:border-green-800">
                            <div className="flex items-start justify-between">
                              <code className="text-sm text-green-900 dark:text-green-100 flex-1">
                                {commands.interactiveCommand}
                              </code>
                              <button
                                onClick={() => copyToClipboard(commands.interactiveCommand)}
                                className="ml-2 text-green-600 hover:text-green-800 dark:text-green-400 dark:hover:text-green-200"
                              >
                                <ClipboardDocumentIcon className="h-4 w-4" />
                              </button>
                            </div>
                          </div>
                        </div>

                        {/* Manual download */}
                        <details className="group">
                          <summary className="text-xs font-medium text-gray-700 dark:text-gray-300 cursor-pointer hover:text-gray-900 dark:hover:text-gray-100">
                            📥 Manual Download & Install
                          </summary>
                          <div className="mt-2 bg-gray-50 dark:bg-gray-800 rounded-md p-3 border border-gray-200 dark:border-gray-700">
                            <div className="flex items-start justify-between">
                              <code className="text-sm text-gray-900 dark:text-white flex-1 whitespace-pre-line">
                                {commands.manualCommand}
                              </code>
                              <button
                                onClick={() => copyToClipboard(commands.manualCommand)}
                                className="ml-2 text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
                              >
                                <ClipboardDocumentIcon className="h-4 w-4" />
                              </button>
                            </div>
                          </div>
                        </details>
                      </div>
                    );
                  })()}
                </div>
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
                    Sensor Name *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="sensor-dc01"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    IP Address *
                  </label>
                  <input
                    type="text"
                    value={formData.ipAddress}
                    onChange={(e) => setFormData({...formData, ipAddress: e.target.value})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="192.168.1.100"
                  />
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    The sensor host must have this IP address assigned to one of its interfaces
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Tags
                  </label>
                  <div className="mt-2 space-y-2">
                    {availableTags.map((tag) => (
                      <label key={tag} className="inline-flex items-center">
                        <input
                          type="checkbox"
                          checked={formData.tags.includes(tag)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setFormData({...formData, tags: [...formData.tags, tag]});
                            } else {
                              setFormData({...formData, tags: formData.tags.filter(t => t !== tag)});
                            }
                          }}
                          className="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
                        />
                        <span className="ml-2 text-sm text-gray-700 dark:text-gray-300">{tag}</span>
                      </label>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Profile
                  </label>
                  <select
                    value={formData.profile}
                    onChange={(e) => setFormData({...formData, profile: e.target.value})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                  >
                    {availableProfiles.map((profile) => (
                      <option key={profile.value} value={profile.value}>
                        {profile.label} - {profile.description}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                    rows={3}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="Optional description for this sensor"
                  />
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
                  onClick={handleAddSensor}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  Generate Key
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Admin Settings Modal */}
      {showAdminSettings && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white dark:bg-gray-800">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                Admin Settings
              </h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Key Expiration (minutes)
                  </label>
                  <input
                    type="number"
                    value={adminSettings.keyExpirationMinutes}
                    onChange={(e) => setAdminSettings({...adminSettings, keyExpirationMinutes: parseInt(e.target.value)})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    min="5"
                    max="1440"
                  />
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    How long registration keys remain valid (5-1440 minutes)
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Max Pending Sensors
                  </label>
                  <input
                    type="number"
                    value={adminSettings.maxPendingSensors}
                    onChange={(e) => setAdminSettings({...adminSettings, maxPendingSensors: parseInt(e.target.value)})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    min="1"
                    max="1000"
                  />
                </div>

                <div className="flex items-center">
                  <input
                    type="checkbox"
                    checked={adminSettings.requireIpValidation}
                    onChange={(e) => setAdminSettings({...adminSettings, requireIpValidation: e.target.checked})}
                    className="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
                  />
                  <label className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                    Require IP Address Validation
                  </label>
                </div>
              </div>
              <div className="flex justify-end space-x-2 mt-6">
                <button
                  onClick={() => setShowAdminSettings(false)}
                  className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
                >
                  Cancel
                </button>
                <button
                  onClick={() => setShowAdminSettings(false)}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  Save Settings
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SensorRegistrationPage;
