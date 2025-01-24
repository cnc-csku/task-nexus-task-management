import axios from '@/lib/axios/axios.config';
import { NextRequest, NextResponse } from 'next/server';
import { SetupStatus } from './interfaces/Setup';

export async function middleware(req: NextRequest) {
    const pathName = req.nextUrl.pathname;

    // Skip the middleware for static files
    if (pathName.startsWith('/_next/static') || pathName === '/favicon.ico' || pathName.startsWith('/api/auth')) {
        return NextResponse.next();
    }

    try {
        // Call the API to check setup status
        const setupResponse = await axios.get<SetupStatus>('/setup/v1');

        if (setupResponse.status === 200) {
            const setupStatus = setupResponse.data;

            if (!pathName.startsWith('/setup') && setupStatus.isSetupOwner && setupStatus.isSetupWorkspace) {
                return NextResponse.next();
            }

            if (pathName.startsWith('/setup') && setupStatus.isSetupOwner && setupStatus.isSetupWorkspace) {
                return NextResponse.redirect(new URL('/', req.url));
            }

            if (!pathName.startsWith('/setup') || setupStatus.isSetupOwner && setupStatus.isSetupWorkspace) {
                return NextResponse.redirect(new URL('/setup', req.url));
            } else {
                return NextResponse.next();
            }
        }
    } catch (error) {
        console.error('Error checking setup status:', error);
    }

    // Redirect to `/setup` if the conditions are not met or an error occurs
    return NextResponse.redirect(new URL('/setup', req.url));
}

// Apply the middleware to all routes
export const config = {
    matcher: '/:path*', // Matches all paths
};
