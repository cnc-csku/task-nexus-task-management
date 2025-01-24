import { UserLoginRequest, UserLoginResponse } from '@/interfaces/User';
import axios from '@/lib/axios/axios.config';
import { AuthOptions } from 'next-auth';
import Credentials from 'next-auth/providers/credentials';

export const authOptions: AuthOptions = {
  debug: true,
  providers: [
    Credentials({
      credentials: {
        email: {},
        password: {},
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null
        }

        try {
          const payload: UserLoginRequest = {
            email: credentials.email,
            password: credentials.password,
          }

          const user = await axios.post<UserLoginResponse>('/auth/v1/login', payload);

          const response = user.data;

          return {
            id: response.id,
            email: response.email,
            fullName: response.fullName,
            displayName: response.displayName,
            profileUrl: response.profileUrl,
            token: response.token,
            tokenExpireAt: new Date(response.tokenExpireAt),
            createdAt: new Date(response.createdAt),
            updatedAt: new Date(response.updatedAt),
          }

        } catch (error) {
          console.log('Auth Error', error);
          return null
        }
      },
    }),
  ],
  session: {
    strategy: 'jwt',
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.id = user.id;
        token.email = user.email;
        token.fullName = user.fullName;
        token.displayName = user.displayName;
        token.profileUrl = user.profileUrl;
        token.token = user.token;
        token.tokenExpireAt = user.tokenExpireAt;
        token.createdAt = user.createdAt;
        token.updatedAt = user.updatedAt;
      }

      return token;
    },
    async session({ session, token }) {
      if (token) {
        session.user = token;
      }

      return session;
    }
  },
  pages: {
    signIn: '/auth/login',
  },
};
