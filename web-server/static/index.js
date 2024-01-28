const host = 'http://localhost:80';
const jwtCookieName = 'auth_jwt';

function app() {
  const state = {
    isAuthenticated: false,
    data: []
  }

  const body = document.querySelector('body');
  
  checkIsAuthenticated();
  render();

  function render() {
    body.innerHTML = '';
    if (state.isAuthenticated) {
      renderMenu();
      renderSessionsData();
      return;
    }
    renderAuthForm()
  }

  async function renderAuthForm() {
    const container = document.createElement('div');
    container.classList = 'container';
    form = document.createElement('form');
    form.classList = 'container-sm border auth-form has-validation auth-form';
    form.innerHTML = `<div class="mb-3" >
      <label for="login" class="form-label">Login</label>
      <input type="login" class="form-control" id="login" aria-describedby="loginHelp">
      <div class="invalid-feedback">
        Please enter a message in the textarea.
      </div>
    </div>
    <div class="mb-3">
      <label for="password" class="form-label">Password</label>
      <input type="password" class="form-control" id="password">
    </div>
    <div class="text-danger error-info"></div>
    <button type="submit" class="btn btn-primary">Submit</button>`
    form.onsubmit = (e) => { handleAuth(e) };

    container.appendChild(form);
    body.appendChild(container);

    async function handleAuth(e) {
      e.preventDefault();
      const username = document.getElementById("login").value;
      const password = document.getElementById("password").value;
      if (username === "" || password === "") {
        const errorDiv = document.querySelector('.error-info');
        errorDiv.innerHTML = `Login and password should not be empty`;
        return;
      }
      await login(username, password);
      if (state.isAuthenticated) {
        render();
      }
    }
  
    async function login(login, password) {
      const url = `${host}/login`;
      const xhr = new XMLHttpRequest();
      const data = JSON.stringify({"login": login, "password": password});
      
      return new Promise(function(resolve, reject) {
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onload = function () {
          if (this.status >= 200 && this.status < 300) {
            const jwt = xhr.responseText;
            setJWTTokenToCookie(jwt);
            state.isAuthenticated = true;
            resolve();
          } else {
            const errorDiv = document.querySelector('.error-info');
            errorDiv.innerHTML = `Invalid login or password. Please try again`;
            reject()
          }
        };
        xhr.send(data);
      })
    }
  
    function setJWTTokenToCookie(jwt) {
      var currentDate = new Date();
      currentDate.setTime(currentDate.getTime() + (24 * 60 * 60 * 1000));
      var expires = "expires=" + currentDate.toUTCString();
      document.cookie = `auth_jwt=${jwt}; ` + expires + "; path=/";
    }
  }

  function renderMenu() {
    const navElement = document.createElement('nav');
    navElement.classList = 'navbar navbar-expand-lg bg-body-tertiary';
    navElement.innerHTML = `<div class="container-fluid"><div class="collapse navbar-collapse" id="navbarScroll"> <ul class="navbar-nav me-auto my-2 my-lg-0 navbar-nav-scroll" style="--bs-scroll-height: 100px;"></ul></div></div>`
    const logoutButton = document.createElement('button');
    logoutButton.classList = 'btn btn-primary';
    logoutButton.addEventListener('click', logout);
    logoutButton.innerText = 'Logout';
    document.querySelector('body').appendChild(navElement);
    document.querySelector('.navbar-collapse').appendChild(logoutButton);
    
    function logout() {
      document.cookie = `${jwtCookieName}=; Max-Age=-99999999;`;
      state.isAuthenticated = false;
      render();
    }
  }

  async function renderSessionsData() {
    const container = document.createElement('div');
    container.classList = 'container';
    const h1 = document.createElement('h1');
    h1.innerHTML = 'Данные о сессиях пользователей';
    
    const jwt = getCookie(jwtCookieName);
    const sessions = await requestSessions(jwt);
    const table = createSessionsTable(sessions);

    container.appendChild(h1);
    container.appendChild(table);

    body.appendChild(container);

    async function requestSessions(jwt) {
      const url = `${host}/data`;
      const options = {
        method: 'GET',
        headers: {
          'Authorization': jwt,
        }
      };
      return await fetch(url, options).then((response) => response.json());
    }

    function createSessionsTable(sessions) {
      const table = document.createElement('table');
      table.classList = 'table';
      
      const thead = document.createElement('thead');
      thead.innerHTML = `<tr><th scope="col">#</th><th scope="col">Session started at</th><th scope="col">Session ended at</th><th scope="col">Client IP</th><th scope="col">SQL Queries</th></tr>`
      const tbody = document.createElement('tbody');
      sessions.forEach((session, index) => {
        const tr = document.createElement('tr');
        tr.appendChild(createTd(index + 1));
        tr.appendChild(createTd(session.started_at));
        tr.appendChild(createTd(session.started_at));
        tr.appendChild(createTd(session.client_ip));
        tr.appendChild(createQueriesTd(session.sql_queries));
        tbody.appendChild(tr);
      });

      table.appendChild(thead);
      table.appendChild(tbody);
      return table;
    }

    function createTd(data) {
      const td = document.createElement('td');
      td.innerHTML = data;
      return td;
    }

    function createQueriesTd(queries) {
      const queriesTd = document.createElement('td');
      queries.forEach((query) => {
        const spanText = query.is_success ? 'success' : 'failure';
        const spanClass = query.is_success ? 'text-success' : 'text-danger';
        queriesTd.innerHTML += `<code>${query.query}</code> - <span class=${spanClass}>${spanText}</span></br>`
      })
      return queriesTd;
    }
  }

  function checkIsAuthenticated() {
    const jwt = getCookie(jwtCookieName);
    if (jwt === '' || jwt === null || jwt === undefined) {
      state.isAuthenticated = false;
      return;
    }
    state.isAuthenticated = true;
  }

  function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i < ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0)==' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
  }    
}

app();