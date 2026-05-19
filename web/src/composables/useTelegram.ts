export function useTelegram() {
  const tg = window.Telegram?.WebApp

  function hapticImpact(style: 'light' | 'medium' | 'heavy' = 'light') {
    tg?.HapticFeedback?.impactOccurred(style)
  }

  function hapticSelection() {
    tg?.HapticFeedback?.selectionChanged()
  }

  function hapticNotification(type: 'success' | 'error' | 'warning') {
    tg?.HapticFeedback?.notificationOccurred(type)
  }

  function showMainButton(text: string, callback: () => void) {
    if (!tg?.MainButton) return
    tg.MainButton.offClick(callback)
    tg.MainButton.setText(text)
    tg.MainButton.onClick(callback)
    tg.MainButton.show()
  }

  function hideMainButton() {
    tg?.MainButton?.hide()
  }

  function showBackButton(callback: () => void) {
    if (!tg?.BackButton) return
    tg.BackButton.offClick(callback)
    tg.BackButton.onClick(callback)
    tg.BackButton.show()
  }

  function hideBackButton() {
    tg?.BackButton?.hide()
  }

  // setClosingGuard toggles Telegram's native "Закрыть приложение?" confirmation
  // for the WebView's close button. Enable while there's in-progress user work
  // (active workout, dirty form, mid-onboarding) and disable on unmount.
  function setClosingGuard(enabled: boolean) {
    if (!tg) return
    if (enabled) {
      tg.enableClosingConfirmation?.()
    } else {
      tg.disableClosingConfirmation?.()
    }
  }

  return {
    tg,
    hapticImpact,
    hapticSelection,
    hapticNotification,
    showMainButton,
    hideMainButton,
    showBackButton,
    hideBackButton,
    setClosingGuard,
  }
}
