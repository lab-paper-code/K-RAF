const request = async (operation, url, params, headers, username, password, credentials = 'include') => {
  let method = operation.toUpperCase();
  let _url = 'http://155.230.36.27:31200' + url;

  if (method === 'GET') {
    _url += new URLSearchParams(params).toString();
  }

  const base64Credentials = btoa(username + ":" + password);

  let options = {
    method: method,
    headers: {
      "Content-Type": 'application/json',
      "Authorization": `Basic ${base64Credentials}`,
      ...headers
    },
    credentials: 'include',
  };
  
  if (method !== 'GET') {
    options['body'] = JSON.stringify(params);
  }
  const response = await fetch(_url, options);
  if (response.status === 204 || response.status === 205) {
    // No content or Reset content, no need to parse JSON
    return null;
  } else if (response.status >= 200 && response.status < 300) {
    return await response.json();
  } else {
    const errorData = await response.json();
    throw new Error(JSON.stringify(errorData));
  }
};

export default request;
