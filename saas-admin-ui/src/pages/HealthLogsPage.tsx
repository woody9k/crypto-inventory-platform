import { useEffect, useState } from 'react';
import api from '../services/api';

export default function HealthLogsPage() {
  const [health, setHealth] = useState<any>(null);
  const [logs, setLogs] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        const [h, l] = await Promise.all([
          api.get('/admin/monitoring/health'),
          api.get('/admin/monitoring/logs'),
        ]);
        setHealth(h.data);
        setLogs(l.data);
      } catch (e: any) {
        setError(e?.response?.data?.error || 'Failed to load monitoring data');
      }
    };
    load();
  }, []);

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-semibold">Platform Monitoring</h1>
      {error && <div className="text-red-600">{error}</div>}
      <div>
        <h2 className="font-medium mb-2">Health</h2>
        <pre className="bg-gray-100 p-3 rounded text-sm overflow-auto">{JSON.stringify(health, null, 2)}</pre>
      </div>
      <div>
        <h2 className="font-medium mb-2">Logs</h2>
        <pre className="bg-gray-100 p-3 rounded text-sm overflow-auto">{JSON.stringify(logs, null, 2)}</pre>
      </div>
    </div>
  );
}


