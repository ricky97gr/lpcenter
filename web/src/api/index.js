import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  login(data) {
    return api.post('/auth/login', data)
  },
  logout() {
    return api.post('/auth/logout')
  }
}

export const userAPI = {
  getAll(params = {}) {
    return api.get('/users', { params })
  },
  getById(id) {
    return api.get(`/users/${id}`)
  },
  getCurrent() {
    return api.get('/users/me')
  },
  create(data) {
    return api.post('/users', data)
  },
  update(id, data) {
    return api.put(`/users/${id}`, data)
  },
  delete(id) {
    return api.delete(`/users/${id}`)
  },
  resetPassword(id, data) {
    return api.post(`/users/${id}/reset-password`, data)
  },
  changePassword(data) {
    return api.post('/users/change-password', data)
  }
}

export const dashboardAPI = {
  getStats() {
    return api.get('/dashboard/stats')
  },
  getRecentLicenses() {
    return api.get('/dashboard/recent-licenses')
  },
  getRecentPlugins() {
    return api.get('/dashboard/recent-plugins')
  }
}

export const licenseTypeAPI = {
  getAll(params = {}) {
    return api.get('/license-types', { params })
  },
  getById(id) {
    return api.get(`/license-types/${id}`)
  },
  create(data) {
    return api.post('/license-types', data)
  },
  update(id, data) {
    return api.put(`/license-types/${id}`, data)
  },
  delete(id) {
    return api.delete(`/license-types/${id}`)
  }
}

export const productAPI = {
  getAll(params = {}) {
    return api.get('/products', { params })
  },
  getById(id) {
    return api.get(`/products/${id}`)
  },
  create(data) {
    return api.post('/products', data)
  },
  update(id, data) {
    return api.put(`/products/${id}`, data)
  },
  delete(id) {
    return api.delete(`/products/${id}`)
  }
}

export const publicAPI = {
  createPlugin(data) {
    return api.post('/public/plugins', data)
  },
  uploadPlugin(data) {
    return api.post('/public/plugins/upload', data, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  getProducts(params = {}) {
    return api.get('/public/products', { params })
  },
  getLicenseTypes(params = {}) {
    return api.get('/public/license-types', { params })
  }
}

export const pluginAPI = {
  getAll(params = {}) {
    return api.get('/plugins', { params })
  },
  getById(id) {
    return api.get(`/plugins/${id}`)
  },
  create(data) {
    return api.post('/plugins', data)
  },
  update(id, data) {
    return api.put(`/plugins/${id}`, data)
  },
  updateStatus(id, data) {
    return api.put(`/plugins/${id}/status`, data)
  },
  sign(id) {
    return api.put(`/plugins/${id}/sign`)
  },
  publish(id) {
    return api.put(`/plugins/${id}/publish`)
  },
  disable(id) {
    return api.put(`/plugins/${id}/disable`)
  },
  delete(id) {
    return api.delete(`/plugins/${id}`)
  }
}

export const licenseAPI = {
  getAll(params = {}) {
    return api.get('/licenses', { params })
  },
  getById(id) {
    return api.get(`/licenses/${id}`)
  },
  create(data) {
    return api.post('/licenses', data)
  },
  update(id, data) {
    return api.put(`/licenses/${id}`, data)
  },
  delete(id) {
    return api.delete(`/licenses/${id}`)
  },
  approve(id) {
    return api.put(`/licenses/${id}/approve`)
  },
  reject(id) {
    return api.put(`/licenses/${id}/reject`)
  },
  download(id) {
    return api.get(`/licenses/${id}/download`, { responseType: 'blob' })
  }
}

export default api