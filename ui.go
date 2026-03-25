package main

import "fmt"

const CalmStyles = `
<style>
    :root { --bg: #050505; --panel: #111111; --border: #1f1f1f; --text: #e0e0e0; --accent: #3498db; --danger: #e74c3c; }
    body { background: var(--bg); color: var(--text); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 0; display: flex; height: 100vh; overflow: hidden; }
    
    /* Авторизация */
    .auth-container { display: flex; align-items: center; justify-content: center; width: 100vw; height: 100vh; }
    .card { background: var(--panel); padding: 40px; border-radius: 24px; border: 1px solid var(--border); width: 100%; max-width: 360px; box-shadow: 0 10px 40px rgba(0,0,0,0.5); text-align: center; }
    .avatar-upload { width: 80px; height: 80px; background: #1a1a1a; border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; border: 1px dashed #444; cursor: pointer; color: #666; transition: 0.3s; }
    .avatar-upload:hover { border-color: var(--accent); color: var(--accent); }
    input { width: 100%; background: #1a1a1a; border: 1px solid var(--border); color: #fff; padding: 16px; border-radius: 14px; margin-bottom: 12px; box-sizing: border-box; outline: none; transition: 0.3s; }
    input:focus { border-color: var(--accent); }
    .btn { background: var(--accent); color: #fff; border: none; width: 100%; padding: 16px; border-radius: 14px; font-weight: 600; font-size: 16px; cursor: pointer; transition: 0.3s; }
    .btn:hover { opacity: 0.8; transform: translateY(-1px); }

    /* Интерфейс Чатов */
    .sidebar { width: 350px; background: var(--panel); border-right: 1px solid var(--border); display: flex; flex-direction: column; position: relative; }
    .top-bar { padding: 15px; display: flex; gap: 10px; align-items: center; border-bottom: 1px solid var(--border); }
    .search { flex-grow: 1; margin: 0; padding: 12px 15px; border-radius: 10px; background: #1a1a1a; }
    .icon-btn { background: none; border: none; font-size: 20px; cursor: pointer; color: var(--accent); }
    .chat-area { flex-grow: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; background: radial-gradient(circle at center, #0a0a0a, var(--bg)); }
    .welcome-badge { background: rgba(52, 152, 219, 0.1); color: var(--accent); padding: 10px 20px; border-radius: 20px; border: 1px solid rgba(52, 152, 219, 0.2); }
    
    /* Плюсик и Папки */
    .fab { position: absolute; bottom: 80px; right: 20px; width: 56px; height: 56px; background: var(--accent); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 28px; color: #fff; cursor: pointer; text-decoration: none; box-shadow: 0 4px 15px rgba(52, 152, 219, 0.3); transition: 0.3s; }
    .fab:hover { transform: scale(1.05); }
    .nav-bottom { padding: 15px; border-top: 1px solid var(--border); display: flex; justify-content: space-around; background: var(--panel); }
    
    /* Настройки */
    .folder { background: #1a1a1a; padding: 16px; margin-bottom: 10px; border-radius: 14px; display: flex; align-items: center; gap: 15px; cursor: pointer; transition: 0.2s; border: 1px solid transparent; }
    .folder:hover { border-color: var(--border); background: #222; }

    @media (max-width: 768px) { .sidebar { width: 100%; } .chat-area { display: none; } }
</style>`

func GetAuthPage() string {
 return `<!DOCTYPE html><html lang="ru"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>ZIG Login</title>` + CalmStyles + `</head><body>
    <div class="auth-container">
        <div class="card">
            <h2 style="margin-top:0; font-weight:800; color:var(--accent);">ZIG GLOBAL</h2>
            <div class="avatar-upload">📷 Галерея</div>
            <form action="/register" method="POST">
                <input type="text" name="name" placeholder="Ваше Имя" required>
                <input type="text" name="username" placeholder="Username" required>
                <input type="email" name="email" placeholder="E-mail (Gmail)" required>
                <button type="submit" class="btn">Получить код</button>
            </form>
        </div>
    </div></body></html>`
}

func GetVerifyPage(email string) string {
 return `<!DOCTYPE html><html lang="ru"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">` + CalmStyles + `</head><body>
    <div class="auth-container">
        <div class="card">
            <h2>Введите код</h2>
            <p style="color:#888; font-size:14px;">Отправлен на ` + email + `</p>
            <form action="/verify" method="POST">
                <input type="hidden" name="email" value="` + email + `">
                <input type="text" name="code" placeholder="000000" style="letter-spacing:8px; text-align:center; font-size:24px;" maxlength="6" required>
                <input type="password" name="password" placeholder="Облачный пароль (если есть)">
                <button type="submit" class="btn">Войти</button>
            </form>
        </div>
    </div></body></html>`
}

func GetAppLayout() string {
 return `<!DOCTYPE html><html lang="ru"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">` + CalmStyles + `</head><body>
    <div class="sidebar">
        <div class="top-bar">
            <input type="text" class="search" placeholder="Поиск (юзер или чат)...">
            <button class="icon-btn">🔍</button>
        </div>
        <div style="flex-grow:1; padding:20px; color:#555; text-align:center; margin-top:50px;">
            Ничего не найдено
        </div>
        
        <a href="#" class="fab" title="Создать канал/чат">+</a>
        
        <div class="nav-bottom">
            <button class="icon-btn" style="color:var(--text);">💬 Чаты</button>
            <button class="icon-btn" style="color:var(--text);" onclick="document.getElementById('settings').style.display='flex'">⚙️ Настройки</button>
        </div>
    </div>

    <div id="settings" style="display:none; position:absolute; top:0; left:0; width:100%; height:100%; background:var(--bg); z-index:100; flex-direction:column;">
        <div class="top-bar"><button class="icon-btn" onclick="document.getElementById('settings').style.display='none'">🔙 Назад</button> <h3 style="margin:0;">Настройки</h3></div>
        <div style="padding:20px; max-width:600px; margin:0 auto; width:100%; box-sizing:border-box;">
            <div class="folder">🎨 <span>Стиль текста, дизайн, цвета чатов</span></div>
            <div class="folder">🌍 <span>Язык (UA, BE, EN, RU)</span></div>
            <div class="folder">🔐 <span>Облачный пароль и Подсказка</span></div>
            <div class="folder" style="border-color:var(--accent);">👑 <span style="color:var(--accent);">ZIG PRO (Ввести промокод)</span></div>
            <div class="folder" style="color:var(--danger);" onclick="if(confirm('Удалить аккаунт?')){if(confirm('Точно? Блокировка почты 24 часа!')){window.location='/delete?email=user'}}">❌ <span>Удалить аккаунт</span></div>
        </div>
    </div>

    <div class="chat-area">
        <div class="welcome-badge">Вы зарегистрировались на ZIG Global</div>
    </div>
</body></html>`
}

func GetAdminPanel() string {
 return `<!DOCTYPE html><html lang="ru"><head><meta charset="UTF-8">` + CalmStyles + `</head><body>
    <div style="padding:40px; width:100%;">
        <h1 style="color:#FFD700;">God Mode: Панель Создателя</h1>
        <div class="card" style="text-align:left; max-width:600px;">
            <p>Статус: <b>Онлайн</b></p>
            <p>Количество пользователей: <b style="color:var(--accent); font-size:24px;">1</b></p>
            <hr style="border:0; border-top:1px solid var(--border); margin:20px 0;">
            <button class="btn" onclick="alert('Промокод ZIG_PRO_2026 активен')">Управление промокодами</button>
        </div>
        </div></body></html>`
}