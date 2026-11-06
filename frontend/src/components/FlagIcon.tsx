interface Props {
  countryCode: string
  size?: number
}

export default function FlagIcon({ countryCode, size = 14 }: Props) {
  if (!countryCode) return null

  const code = countryCode.toUpperCase()
  const flag = code
    .split('')
    .map(c => String.fromCodePoint(0x1F1E6 + c.charCodeAt(0) - 65))
    .join('')

  return (
    <span style={{ fontSize: size, lineHeight: 1 }} title={code}>
      {flag}
    </span>
  )
}
