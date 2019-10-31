import * as grpc from 'grpc';
import { Observable } from 'rxjs';
/** Properties of a SendPkg. */
export interface SendPkg {

    /** SendPkg data */
    data?: (Uint8Array|null);

    /** SendPkg type */
    type?: (number|null);

    /** SendPkg medata */
    medata?: ({ [k: string]: Uint8Array }|null);
}

/** Properties of a SendRequest. */
export interface SendRequest {

    /** SendRequest to */
    to?: (string|null);

    /** SendRequest token */
    token?: (string|null);

    /** SendRequest pkg */
    pkg?: (ISendPkg|null);
}

/** Properties of a RecvResponse. */
export interface RecvResponse {

    /** RecvResponse from */
    from?: (string|null);

    /** RecvResponse pkg */
    pkg?: (ISendPkg|null);
}

/** Properties of a StatusResponse. */
export interface StatusResponse {

    /** StatusResponse Status */
    Status?: (boolean|null);
}

/** Properties of a RecvRequest. */
export interface RecvRequest {

    /** RecvRequest token */
    token?: (string|null);
}

/** Properties of a UserRequest. */
export interface UserRequest {

    /** UserRequest name */
    name?: (string|null);

    /** UserRequest pwd */
    pwd?: (string|null);
}

/** Properties of a LoginResponse. */
export interface LoginResponse {

    /** LoginResponse token */
    token?: (string|null);

    /** LoginResponse status */
    status?: (boolean|null);

    /** LoginResponse expTime */
    expTime?: (number|Long|null);
}

/** Properties of a LogoutRequest. */
export interface LogoutRequest {

    /** LogoutRequest token */
    token?: (string|null);
}

/** Constructs a new Chat service. */
export interface Chat {

    /**
     * Calls Login.
     * @param request UserRequest message or plain object
     *  * @param metadata Optional metadata
     * @returns Promise
     */
    login(request: IUserRequest, metadata?: grpc.Metadata): Observable<LoginResponse>;

    /**
     * Calls Register.
     * @param request UserRequest message or plain object
     *  * @param metadata Optional metadata
     * @returns Promise
     */
    register(request: IUserRequest, metadata?: grpc.Metadata): Observable<StatusResponse>;

    /**
     * Calls Logout.
     * @param request LogoutRequest message or plain object
     *  * @param metadata Optional metadata
     * @returns Promise
     */
    logout(request: ILogoutRequest, metadata?: grpc.Metadata): Observable<StatusResponse>;

    /**
     * Calls Recv.
     * @param request RecvRequest message or plain object
     *  * @param metadata Optional metadata
     * @returns Promise
     */
    recv(request: IRecvRequest, metadata?: grpc.Metadata): Observable<RecvResponse>;

    /**
     * Calls Send.
     * @param request SendRequest message or plain object
     *  * @param metadata Optional metadata
     * @returns Promise
     */
    send(request: ISendRequest, metadata?: grpc.Metadata): Observable<StatusResponse>;
}
