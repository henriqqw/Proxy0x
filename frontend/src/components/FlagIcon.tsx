import type { ComponentType, SVGProps } from 'react'
import * as Flags from 'country-flag-icons/react/3x2'

type Flag = ComponentType<SVGProps<SVGSVGElement> & { title?: string }>
const flags = Flags as unknown as Record<string, Flag>

interface Props {
  countryCode: string
  size?: number
}

export default function FlagIcon({ countryCode, size = 14 }: Props) {
  const code = countryCode?.toUpperCase()
  const Flag = code ? flags[code] : undefined
  if (!Flag) return null

  return (
    <Flag
      title={code}
      width={size}
      height={Math.round(size * 2 / 3)}
      className="inline-block shrink-0 rounded-[1px]"
    />
  )
}
