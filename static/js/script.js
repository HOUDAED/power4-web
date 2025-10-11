
function selectLevel(level) {
    const buttons = document.querySelectorAll('.level-btn');
    buttons.forEach(btn => btn.classList.remove('selected'));
    document.querySelector('.' + level).classList.add('selected');
    document.getElementById('difficultyInput').value = level;
    setTimeout(() => {
    document.getElementById('startForm').submit();
    }, 200);
}