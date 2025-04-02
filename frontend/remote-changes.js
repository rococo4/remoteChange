document.addEventListener('DOMContentLoaded', () => {
	// При загрузке страницы вызываем /teams/user/role
	fetch('/teams/user/role')
		.then(res => res.json())
		.then(data => {
			if(data.role !== 'admin') {
				document.getElementById('admin-panel').style.display = 'none';
			}
		})
		.catch(err => console.error(err));

	// Обработчик нажатия на кнопку Team
	document.getElementById('teamBtn').addEventListener('click', () => {
		fetch('/teams/user/team')
			.then(res => res.json())
			.then(teamData => {
				if(!teamData.teamId) {
					document.getElementById('team-info').innerText = 'нет команды';
					return;
				}
				document.getElementById('team-name').innerText = teamData.teamName;
				fetch(`/configs/team/${teamData.teamId}`)
					.then(res => res.json())
					.then(configsData => {
						// ...existing rendering code для отображения конфигов команды...
					})
					.catch(err => console.error(err));
			})
			.catch(err => console.error(err));
	});

	// Обработчик нажатия кнопки Watch
	document.getElementById('watchBtn').addEventListener('click', () => {
		const configId = getSelectedConfigId();
		Promise.all([
			fetch(`/configs/${configId}`).then(res => res.json()),
			fetch(`/configs/${configId}/changes`).then(res => res.json())
		])
			.then(([config, changes]) => {
				// ...existing display code для отображения конфига и истории...
			})
			.catch(err => console.error(err));
	});

	// Обработчик нажатия кнопки Edit
	document.getElementById('editBtn').addEventListener('click', () => {
		const configId = getSelectedConfigId();
		fetch(`/configs/${configId}`)
			.then(res => res.json())
			.then(config => {
				document.getElementById('config-editor').value = config.content;
			})
			.catch(err => console.error(err));
	});

	// Обработчик сохранения изменений (Edit->Save)
	document.getElementById('saveBtn').addEventListener('click', () => {
		const configId = getSelectedConfigId();
		const newContent = document.getElementById('config-editor').value;
		fetch(`/configs/${configId}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ content: newContent })
		})
			.then(res => {
				if(res.ok) {
					// ...existing code для уведомления об успешном обновлении...
				}
			})
			.catch(err => console.error(err));
	});

	// Обработчик кнопки Rollback
	document.getElementById('rollbackBtn').addEventListener('click', () => {
		const configId = getSelectedConfigId();
		fetch(`/configs/${configId}/rollback`, {
			method: 'POST'
		})
			.then(res => {
				if(res.ok) {
					// ...existing code для уведомления о успешном откате...
				}
			})
			.catch(err => console.error(err));
	});
});

// Вспомогательная функция для получения выбранного id конфига
function getSelectedConfigId() {
	// ...existing code для получения id...
	return document.getElementById('config-id').value;
}
