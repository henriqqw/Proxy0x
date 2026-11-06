interface Props {
  countryCode: string
  size?: number
}

export default function FlagIcon({ countryCode, size = 14 }: Props) {
  if (!countryCode) return null

  return (
    <img
      src={`https://flagcdn.com/w20/${countryCode.toLowerCase()}.png`}
      alt={countryCode.toUpperCase()}
      title={countryCode.toUpperCase()}
      width={size}
      height={Math.round(size * 0.75)}
      className="inline-block object-contain"
      style={{ minWidth: size }}
      loading="lazy"
      onError={e => {
        (e.target as HTMLImageElement).style.display = 'none'
      }}
    />
  )
}
