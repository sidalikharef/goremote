import '@testing-library/jest-dom'
import { expect } from 'vitest'

declare global {
  // eslint-disable-next-line no-var
  var expect: typeof import('vitest').expect
}

globalThis.expect = expect
