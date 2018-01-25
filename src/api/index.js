import mocks from './mocks';
import TokenService from '../services/token';

export const serverUrl =
  process.env.REACT_APP_SERVER_URL || 'http://127.0.0.1:8080';

const apiUrl = url => `${serverUrl}${url}`;

function checkStatus(resp) {
  // when server return Unauthorized we need to remove token
  if (resp.status === 401) {
    TokenService.remove();
  }
  if (resp.status < 200 || resp.status >= 300) {
    const error = new Error(resp.statusText);
    error.response = resp;
    throw error;
  }
  return resp;
}

function normalizeError(err) {
  if (typeof err === 'object') {
    // error from server
    if (err.title) {
      return err.title;
    }
    // javascript error
    if (err.message) {
      return err.message;
    }
    // weird object as error, shouldn't really happen
    return JSON.stringify(err);
  }
  if (typeof err === 'string') {
    return err;
  }
  return 'Internal error';
}

function normalizeErrors(err) {
  if (Array.isArray(err)) {
    return err.map(e => normalizeError(e));
  }
  return [normalizeError(err)];
}

function apiCall(url, options = {}) {
  const token = TokenService.get();

  return fetch(apiUrl(url), {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token}`,
    },
  })
    .then(checkStatus)
    .then(resp => {
      // no content
      if (resp.status === 204) {
        return {};
      }

      return resp.json().then(json => {
        if (json.errors) {
          throw json.errors;
        }
        return json.data;
      });
    })
    .catch(err => Promise.reject(normalizeErrors(err)));
}

function me() {
  return apiCall(`/api/me`);
}

function getExperiment(experimentId) {
  return apiCall(`/api/experiments/${experimentId}`);
}

function getAssignments(experimentId) {
  return apiCall(`/api/experiments/${experimentId}/assignments`);
}

function getFilePair(experimentId, pairId) {
  return apiCall(`/api/experiments/${experimentId}/file-pairs/${pairId}`);
}

function putAssigment(id, body) {
  return apiCall(`/api/experiments/1/assignments/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  });
}

export default {
  me,

  getExperiment,
  getAssignments,
  getFilePair,
  putAssigment,
};
