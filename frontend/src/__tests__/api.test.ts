import axios from 'axios';

describe('backend connectivity', () => {
  it('responds to GET /api/go/actors', async () => {
    const base = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';
    const res = await axios.get(`${base}/api/go/actors`);
    expect(res.status).toBe(200);
    // should be an array (even if empty)
    expect(Array.isArray(res.data)).toBe(true);
  });
});
