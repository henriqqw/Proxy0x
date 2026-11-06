interface Props {
  countryCode: string
  size?: number
}

export default function FlagIcon({ countryCode, size = 14 }: Props) {
  if (!countryCode) return null

  const code = countryCode.toLowerCase()

  return (
    <img
      src={`https://flagcdn.com/w40/${code}.png`}
      srcSet={`https://flagcdn.com/w80/${code}.png 2x`}
      alt={countryCode.toUpperCase()}
      title={countryCode.toUpperCase()}
      style={{ width: size, height: Math.round(size * 0.75), objectFit: 'cover', borderRadius: 2 }}
      draggable={false}
      loading="lazy"
    />
  )
}
