// Helpers for the user-side video player.
//
// Admin-uploaded videos live in R2 (public bucket) and the server replaces
// `video_url` with a direct .mp4 URL. External URLs (YouTube etc.) are kept
// as paste-only embeds. The two need different players: <video> for direct
// files, <iframe> for embeds.

export function isDirectVideoUrl(url: string): boolean {
  if (!url) return false
  return /\.mp4(\?|$)/i.test(url)
    || /r2\.cloudflarestorage\.com/.test(url)
    || /pub-[a-z0-9]+\.r2\.dev/.test(url)
}
