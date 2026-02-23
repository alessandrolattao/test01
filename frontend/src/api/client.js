const API_BASE = '/api'

async function request(path, options = {}) {
  const url = `${API_BASE}${path}`
  const headers = { ...options.headers }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  const token = localStorage.getItem('admin_token')
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...options, headers })

  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    const error = new Error(body.error || `Request failed: ${res.status}`)
    error.status = res.status
    throw error
  }

  if (res.status === 204) return null
  return res.json()
}

export function get(path) {
  return request(path)
}

export function post(path, body) {
  if (body instanceof FormData) {
    return request(path, { method: 'POST', body })
  }
  return request(path, { method: 'POST', body: JSON.stringify(body) })
}
