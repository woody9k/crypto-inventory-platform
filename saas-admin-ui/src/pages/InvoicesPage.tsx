import { useEffect, useState } from 'react';
import api from '../services/api';

export default function InvoicesPage() {
  const [invoices, setInvoices] = useState<any[]>([]);
  const [tenantId, setTenantId] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = async () => {
    setLoading(true);
    try {
      const url = tenantId ? `/admin/billing/invoices?tenantId=${tenantId}` : '/admin/billing/invoices';
      const res = await api.get(url);
      setInvoices(res.data.invoices || []);
    } catch (e: any) {
      setError(e?.response?.data?.error || 'Failed to load invoices');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-semibold">Invoices</h1>
      <div className="flex items-center space-x-2">
        <input
          className="border p-1"
          placeholder="Filter by tenantId (optional)"
          value={tenantId}
          onChange={(e) => setTenantId(e.target.value)}
        />
        <button className="bg-gray-800 text-white px-3 py-1 rounded" onClick={load}>Reload</button>
      </div>
      {loading ? (
        <div>Loading...</div>
      ) : error ? (
        <div className="text-red-600">{error}</div>
      ) : (
        <div className="overflow-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="text-left">
                <th className="p-2">Invoice ID</th>
                <th className="p-2">Amount</th>
                <th className="p-2">Currency</th>
                <th className="p-2">Status</th>
                <th className="p-2">Issued</th>
                <th className="p-2">Due</th>
                <th className="p-2">Paid</th>
              </tr>
            </thead>
            <tbody>
              {invoices.map((inv, idx) => (
                <tr key={idx} className="border-t">
                  <td className="p-2">{inv.invoice_id}</td>
                  <td className="p-2">{(inv.amount_cents / 100).toFixed(2)}</td>
                  <td className="p-2">{inv.currency}</td>
                  <td className="p-2">{inv.status}</td>
                  <td className="p-2">{inv.issued_at ? new Date(inv.issued_at).toLocaleString() : '-'}</td>
                  <td className="p-2">{inv.due_at ? new Date(inv.due_at).toLocaleString() : '-'}</td>
                  <td className="p-2">{inv.paid_at ? new Date(inv.paid_at).toLocaleString() : '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}


