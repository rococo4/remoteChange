document.addEventListener('DOMContentLoaded', () => {
    checkUserRole();
    document.getElementById('team-btn').addEventListener('click', handleTeamClick);
    document.getElementById('watch-btn').addEventListener('click', handleWatch);
    document.getElementById('edit-btn').addEventListener('click', handleEdit);
    document.getElementById('save-btn').addEventListener('click', handleSave);
    document.getElementById('rollback-btn').addEventListener('click', handleRollback);
});

function checkUserRole() {
    fetch('/teams/user/role')
      .then(resp => resp.json())
      .then(role => {
          if (role !== 'admin') {
              document.getElementById('admin-panel').style.display = 'none';
          }
      })
      .catch(console.error);
}

function handleTeamClick() {
    fetch('/teams/user/team')
      .then(resp => resp.json())
      .then(team => {
          if (!team || !team.id) {
              document.getElementById('team-info').innerText = 'нет команды';
              return;
          }
          document.getElementById('team-name').innerText = team.name;
          fetch(`/configs/team/${team.id}`)
            .then(resp => resp.json())
            .then(configs => {
                const configsDiv = document.getElementById('team-configs');
                configsDiv.innerHTML = '';
                configs.forEach(config => {
                    const btn = document.createElement('button');
                    btn.innerText = `Config ${config.id}`;
                    btn.addEventListener('click', () => loadConfig(config.id));
                    configsDiv.appendChild(btn);
                });
            });
      })
      .catch(console.error);
}

function loadConfig(configId) {
    // Show action buttons and save configId in dataset for each
    document.getElementById('watch-btn').style.display = 'inline';
    document.getElementById('edit-btn').style.display = 'inline';
    document.getElementById('rollback-btn').style.display = 'inline';
    document.getElementById('watch-btn').dataset.configId = configId;
    document.getElementById('edit-btn').dataset.configId = configId;
    document.getElementById('rollback-btn').dataset.configId = configId;
}

function handleWatch(e) {
    const configId = e.target.dataset.configId;
    fetch(`/configs/${configId}`)
      .then(resp => resp.json())
      .then(config => {
          alert('Config content: ' + config.content);
      });
    fetch(`/configs/${configId}/changes`)
      .then(resp => resp.json())
      .then(changes => {
          document.getElementById('config-history').innerText = JSON.stringify(changes);
      });
}

function handleEdit(e) {
    const configId = e.target.dataset.configId;
    fetch(`/configs/${configId}`)
      .then(resp => resp.json())
      .then(config => {
          const textarea = document.getElementById('edit-area');
          textarea.style.display = 'block';
          textarea.value = config.content;
          document.getElementById('save-btn').style.display = 'inline';
      });
}

function handleSave() {
    const configId = document.getElementById('edit-btn').dataset.configId;
    const newContent = document.getElementById('edit-area').value;
    fetch(`/configs/${configId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: newContent })
    })
    .then(resp => resp.json())
    .then(res => {
        alert('Config updated');
    });
}

function handleRollback(e) {
    const configId = e.target.dataset.configId;
    fetch(`/configs/${configId}/rollback`, { method: 'POST' })
      .then(resp => resp.json())
      .then(res => {
          alert('Config rolled back');
      });
}
