import { useEffect, useState } from 'react';
import { tenantsApi } from '../services/tenantsApi';
import { useParams } from 'react-router-dom';

export default function TenantBillingPage() {
  const { id } = useParams();
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [planKey, setPlanKey] = useState('professional');
  const [cancelAtPeriodEnd, setCancelAtPeriodEnd] = useState(false);

  useEffect(() => {
    if (!id) return;
    tenantsApi
      .getTenantBilling(id)
      .then((res) => setData(res.billing))
      .catch((e) => setError(e?.response?.data?.error || 'Failed to load billing'))
      .finally(() => setLoading(false));
  }, [id]);

  const onChangePlan = async () => {
    if (!id) return;
    await tenantsApi.updateTenantBilling(id, { action: 'change_plan', plan_key: planKey });
    const res = await tenantsApi.getTenantBilling(id);
    setData(res.billing);
  };

  const onToggleCancel = async () => {
    if (!id) return;
    await tenantsApi.updateTenantBilling(id, { action: 'cancel', cancel_at_period_end: !cancelAtPeriodEnd });
    setCancelAtPeriodEnd(!cancelAtPeriodEnd);
    const res = await tenantsApi.getTenantBilling(id);
    setData(res.billing);
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="text-red-600">{error}</div>;

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-semibold">Tenant Billing</h1>
      <pre className="bg-gray-100 p-3 rounded text-sm overflow-auto">{JSON.stringify(data, null, 2)}</pre>

      <div className="space-y-2">
        <div className="flex items-center space-x-2">
          <label className="w-32">Plan</label>
          <select className="border p-1" value={planKey} onChange={(e) => setPlanKey(e.target.value)}>
            <option value="free">Free</option>
            <option value="professional">Professional</option>
            <option value="enterprise">Enterprise</option>
          </select>
          <button className="bg-blue-600 text-white px-3 py-1 rounded" onClick={onChangePlan}>
            Change Plan
          </button>
        </div>

        <div className="flex items-center space-x-2">
          <label className="w-32">Cancel at period end</label>
          <input type="checkbox" checked={cancelAtPeriodEnd} onChange={onToggleCancel} />
        </div>
      </div>
    </div>
  );
}


