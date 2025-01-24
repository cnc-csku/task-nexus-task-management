import { EmployeeRole } from '@/enum/Employee';
import { DefaultSession } from 'next-auth';
import { JWT } from 'next-auth/jwt';

declare module 'next-auth' {
  interface Session {
    user?: User;
    expires: Date;
  }

  interface User {
    id: string;
    email: string;
    fullName: string;
    displayName: string;
    profileUrl: string;
    token: string;
    tokenExpireAt: Date;
    createdAt: Date;
    updatedAt: Date;
  }
}

declare module 'next-auth/jwt' {
  /** Returned by the `jwt` callback and `getToken`, when using JWT sessions */
  interface JWT {
    id: string;
    email: string;
    fullName: string;
    displayName: string;
    profileUrl: string;
    token: string;
    tokenExpireAt: Date;
    createdAt: Date;
    updatedAt: Date;
  }
}
