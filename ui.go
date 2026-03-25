package main

import "fmt"

const CalmStyles = `
<style>
    :root { --bg: #050505; --panel: #111111; --border: #1f1f1f; --text: #e0e0e0; --accent: #007AFF; }
    body { background: var(--bg); color: var(--text); font-family: sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    .card { background: var(--panel); padding: 40px; border-radius: 24px; border: 1px solid var(--border); width: 340px; margin: auto; text-align: center; }
    input { width: 100%; background: #1a1a1a; border: 1px solid var(--border); color: #fff; padding: 15px; border-radius: 12px; margin-bottom: 12px; box-sizing: border-box; }
    .btn { background: var(--accent); color: #fff; border: none; width: 100%; padding: 15px; border-radius: 12px; font-weight: bold; cursor: pointer; }
    .sidebar { width: 350px; background: var(--panel); border-right: 1px solid var(--border); display: flex; flex-direction: column; position: relative; }
    .top-bar { padding: 15px; display: flex; gap: 10px; border-bottom: 1px solid var(--border); }
    .search { flex-grow: 1; padding: 10px; border-radius: 10px; background: #1a1a1a; border: none; color: #fff; }
    .fab { position: absolute; bottom: 80px; right: 20px; width: 60px; height: 60px; background: var(--accent); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 30px; color: #fff; cursor: pointer; text-decoration: none; box-shadow: 0 4px 15px rgba(0,122,255,0.4); }
    .folder { background: #1a1a1a; padding: 15px; margin-bottom: 10px; border-radius: 12px; cursor: pointer; border: 1px solid var(--border); text-align: left; }
    .chat-area { flex-grow: 1; display: flex; align-items: center; justify-content: center; }
</style>`

func GetAuthPage() string {
 return `<!DOCTYPE html><html><head>` + CalmStyles + `</head><body><div class="card">
        <h2 style="color:var(--accent);">ZIG GLOBAL</h2>
        <div style="width:80px; height:80px; background:#222; border-radius:50%; margin:0 auto 20px; line-height:80px;">📷</div>
        <form action="/register" method="POST">
            <input type="text" name="name" placeholder="Имя" required>
            <input type="text" name="username" placeholder="Username" required>
            <input type="email" name="email" placeholder="Gmail" required>
            <button class="btn">Получить код</button>
        </form></div></body></html>`
}

func GetVerifyPage(email string) string {
 // Здесь fmt.Sprintf используется по назначению, ошибка уйдет!
 return fmt.Sprintf(`<!DOCTYPE html><html><head>%s</head><body><div class="card">
        <h2>Введите код</h2><p style="color:#666;">Отправлен на %s</p>
        <form action="/verify" method="POST">
            <input type="hidden" name="email" value="%s">
            <input type="text" name="code" placeholder="000000" maxlength="6" style="text-align:center; font-size:24px; letter-spacing:5px;" required>
            <input type="password" name="password" placeholder="Облачный пароль (если есть)">
            <button class="btn">Войти</button>
        </form></div></body></html>`, CalmStyles, email, email)
}

func GetAppLayout() string {
 return `<!DOCTYPE html><html><head>` + CalmStyles + `</head><body>
    <div class="sidebar">
        <div class="top-bar"><input type="text" class="search" placeholder="Поиск..."><button style="background:none;border:none;color:var(--accent);">🔍</button></div>
        <div style="padding:20px; color:#555; text-align:center;">Ничего не найдено</div>
        <a href="#" class="fab">+</a>
        <div style="padding:15px; border-top:1px solid var(--border); display:flex; justify-content:space-around; background:var(--panel);">
            <button onclick="document.getElementById('settings').style.display='block'" style="background:none; border:none; color:var(--accent); cursor:pointer;">⚙️ Настройки</button>
        </div>
    </div>
    <div id="settings" style="display:none; position:absolute; width:100%; height:100%; background:var(--bg); z-index:100; overflow-y:auto;">
        <div style="padding:20px; max-width:600px; margin:auto; text-align:center;">
            <h2 style="color:var(--accent);">Настройки</h2>
            <div class="folder">🎨 Изменить стиль и дизайн</div>
            <div class="folder">🌍 Язык (RU, UA, BE, EN)</div>
            <div class="folder">🔐 Облачный пароль</div>
            <div class="folder" style="border-color:var(--accent);">👑 Ввести промокод ZIG PRO</div>
            <div class="folder" style="color:#e74c3c;" onclick="if(confirm('Точно удалить? Бан почты на 24 часа!')){window.location='/delete?email=user'}">❌ Удалить аккаунт</div>
            <button class="btn" style="margin-top:20px;" onclick="document.getElementById('settings').style.display='none'">Закрыть</button>
        </div>
    </div>

    <div class="chat-area"><div style="background:#1a1a1a; padding:15px 30px; border-radius:20px; color:var(--accent);">Вы зарегистрировались на ZIG Global</div></div>
    </body></html>`
}

func GetAdminPanel(usersCount int) string {
 return fmt.Sprintf(`<!DOCTYPE html><html><head>%s</head><body><div style="padding:50px; text-align:center; width:100%%;">
        <h1 style="color:#FFD700;">God Mode: Панель Создателя</h1>
        <div class="card" style="margin-top:20px;">
            <p>Статус: <b>Онлайн</b></p>
            <p>Всего пользователей в базе: <b style="color:var(--accent); font-size:24px;">%d</b></p>
        </div></div></body></html>`, CalmStyles, usersCount)
}