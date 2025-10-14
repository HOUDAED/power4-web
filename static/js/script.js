
function selectLevel(level) {
    const buttons = document.querySelectorAll('.level-btn');
    buttons.forEach(btn => btn.classList.remove('selected'));
    document.querySelector('.' + level).classList.add('selected');
    document.getElementById('difficultyInput').value = level;
    setTimeout(() => {
    document.getElementById('startForm').submit();
    }, 200);
}

document.addEventListener('DOMContentLoaded', () => {
  const winSound = document.getElementById('winSound');
  const drawSound = document.getElementById('drawSound');

  // Ces variables doivent être injectées dans ton HTML via templating
  const winner = window.gameState?.winner; // true si victoire
  const draw = window.gameState?.draw;     // true si match nul

  if (winner && winSound) {
    winSound.play();
    triggerConfetti();
  } else if (draw && drawSound) {
    drawSound.play();
  }
});

// 🎊 Confettis avec canvas-confetti
function triggerConfetti() {
  if (typeof confetti !== 'undefined') {
    confetti({
      particleCount: 150,
      spread: 70,
      origin: { y: 0.6 },
      colors: ['#00b4d8', '#ffe66d', '#ff4b4b']
    });
  }
}