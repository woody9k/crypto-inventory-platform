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
import { sensorsApi, PendingSensorItem } from '../services/sensorsApi';

interface AdminSettings {
  keyExpirationMinutes: number;
  maxPendingSensors: number;
  requireIpValidation: boolean;
}

const SensorRegistrationPage: React.FC = () => {
  const [pendingSensors, setPendingSensors] = useState<PendingSensorItem[]>([]);
  const [adminSettings, setAdminSettings] = useState<AdminSettings>({
    keyExpirationMinutes: 60,
    maxPendingSensors: 50,
    requireIpValidation: true
  });
  const [showAddSensor, setShowAddSensor] = useState(false);
  const [showAdminSettings, setShowAdminSettings] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Form state for new sensor
  const [formData, setFormData] = useState({
    name: '',
    ipAddress: '',
    tags: [] as string[],
    profile: 'datacenter_host',
    networkInterfaces: [] as string[],
    description: ''
  });

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

  const fetchPending = async () => {
    try {
      setLoading(true);
      setError(null);
      const items = await sensorsApi.listPending();
      setPendingSensors(items);
    } catch (e) {
      setError('Failed to load pending sensors');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPending();
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

  const generateInstallationCommands = (sensor: PendingSensorItem) => {
    const interfaces = sensor.network_interfaces.length > 0 ? sensor.network_interfaces.join(',') : 'eth0';
    const controlPlaneURL = 'https://crypto-inventory.company.com';
    return {
      curlCommand: `curl -sSL https://sensors.crypto-inventory.com/install.sh | bash -s -- --key ${sensor.registration_key} --ip ${sensor.ip_address} --name "${sensor.name}" --profile ${sensor.profile}`,
      interactiveCommand: `curl -sSL ${controlPlaneURL}/scripts/install-sensor.sh | sudo bash -s -- --interactive`,
      manualCommand: `# Download and run manually
curl -O ${controlPlaneURL}/scripts/install-sensor.sh
chmod +x install-sensor.sh
sudo ./install-sensor.sh \
  --key ${sensor.registration_key} \
  --ip ${sensor.ip_address} \
  --name ${sensor.name} \
  --profile ${sensor.profile} \
  --interfaces "${interfaces}" \
  --url ${controlPlaneURL}`,
    };
  };

  const handleAddSensor = async () => {
    if (!formData.name || !formData.ipAddress) {
      alert('Please fill in required fields');
      return;
    }
    try {
      await sensorsApi.createPending({
        name: formData.name,
        ip_address: formData.ipAddress,
        profile: formData.profile,
        network_interfaces: formData.networkInterfaces,
        tags: formData.tags,
      });
      setShowAddSensor(false);
      setFormData({ name: '', ipAddress: '', tags: [], profile: 'datacenter_host', networkInterfaces: [], description: '' });
      fetchPending();
    } catch (e) {
      setError('Failed to create registration key');
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const deletePendingSensor = async (idOrKey: string) => {
    try {
      await sensorsApi.deletePending(idOrKey);
      fetchPending();
    } catch (e) {
      setError('Failed to delete registration key');
    }
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

      {error && <p className="text-sm text-red-600">{error}</p>}

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
                            Expires in: {formatTimeRemaining(sensor.expires_at)}
                          </span>
                        )}
                      </div>
                      <div className="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                        <p className="mr-4">IP: {sensor.ip_address}</p>
                        <p className="mr-4">Profile: {sensor.profile}</p>
                        <p className="mr-4">Interfaces: {sensor.network_interfaces.join(', ')}</p>
                        <p>Key: {sensor.registration_key}</p>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={() => copyToClipboard(sensor.registration_key)}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4 mr-1" />
                      Copy Key
                    </button>
                    <button
                      onClick={() => copyToClipboard(generateInstallationCommands(sensor).curlCommand)}
                      className="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    >
                      <ClipboardDocumentIcon className="h-4 w-4 mr-1" />
                      Copy Command
                    </button>
                    <button
                      onClick={() => deletePendingSensor(sensor.registration_key)}
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
                    Network Interfaces (comma separated)
                  </label>
                  <input
                    type="text"
                    value={formData.networkInterfaces.join(',')}
                    onChange={(e) => setFormData({...formData, networkInterfaces: e.target.value.split(',').map(s => s.trim()).filter(Boolean)})}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                    placeholder="eth0,eth1"
                  />
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
