import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

export const pluginAPI = {
  getAll() {
    return api.get('/plugins')
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
  delete(id) {
    return api.delete(`/plugins/${id}`)
  }
}

export const licenseAPI = {
  getAll() {
    return api.get('/licenses')
  },
  getById(id) {
    return api.get(`/licenses/${id}`)
  },
  create(data) {
    return api.post('/licenses', data)
  },
  approve(id) {
    return api.put(`/licenses/${id}/approve`)
  },
  reject(id) {
    return api.put(`/licenses/${id}/reject`)
  },
  getByUser(userId) {
    return api.get(`/licenses/user/${userId}`)
  }
}

export default api
