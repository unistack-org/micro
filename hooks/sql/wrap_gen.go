// Code generated. DO NOT EDIT.

package sql

import "database/sql/driver"

func wrapConn(dc driver.Conn, opts Options) driver.Conn {
	c := &wrapperConn{conn: dc, opts: opts}
	if _, ok := dc.(wrapConn0010_77c7c6cd21c875211bd8f7fd0d4f7cac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_c38f9867bd2650446eff3934abfe08ea); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_03ea25d7f45e22ec940d5495f1a84793); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_79f6c233903bd90e363a740e21dedd28); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_396e4ad5df4f8a9883b7db53e3150de5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_58e530768e77481673ca16ec9a1c4151); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_de85b18ab3bb83467b88d133c472e731); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_cf4b2a19600ddd6a199103b184a7b119); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_cc0840d8293f7ce49712448b9ab56686); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_3dd8e3401bce36381241be68e71a4596); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0009_d275e10eaed0aa246adacf06fbc7cb4d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_eb0574608d6a2bdeeda60227474f2696); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_b8a1712d42513721bb27fbf514e4dc08); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_7f39eb6d8011a9d621a372f56c9dbbb9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_4efd319fba1b2f67427b7665301f5965); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_9741ae76b65cf39203ad2ebb1b042b26); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_58b49dcfceedf1a9157418d32ddc3566); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_5cf158fc22b2fa1ebca261e0519f686c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_5392d6d522ce153424f2d6a56e0b9412); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_a6448007f6fd3c5e82ac727e32dd9421); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_8a4cbbf06de72f30a062ff7d0d17a38f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_b62b5d9a91718c412fc82374a1ecc50d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_0ad7ad486ff238407782ecae921802ff); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_7030500b1e198b3921908510db90b60c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_fcea5f14de6e52ea3dcf885d9bbe6c6c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_77ff06f9507fa8e8907ea056fd8b3ff1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_0af8ed4c453b8ed19240067f1ffc2691); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_92d776275e9ff04873b00304c46181c4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_338d30b06f6f46413259d40fd65e6431); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_9424622cd873b695ec4a68179ae1c0d5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_f98a60ed56dc805bbb076e48569dfc95); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_638eb5dd5a685d42b2b21234d0f60533); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_3cca6f68d8b372d89bd75ae899a89fbb); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_93d6b726be36618017eb3d898d1485a9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_79db6e70a15663317e248e989588499f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_83f7a44dbd5760c38c650a39e055e040); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_d5c1b115a4499cda31e8bc5417a3b20c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_d051dedc6dc0f7f58c5f991183648c05); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_519d23a1000bcef9e0c7ba7f65e2d8f4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_a10c9902556dbe03fe0cb57fa734bae3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_d7069a0ed4e229445f0a98eb0a10e6d7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_31cab3e69b231de68861be7149ba4e29); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_a6e602e3b87f67896f9721167c748010); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_f10e8449922c36b75361a92bcd10a530); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_c5919040d602cb66dcdd424433bb7a3a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_3b2f99cd2da334bd09d06e4def61c4ca); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_5d41cafd5772738ac2731fcc2a5feac6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_656773dbc70d65eacff8fb4b5f9e0aa0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_4793858d8022e719f6eaf66c50c2d3ac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_d45f7a4c9249a530fdbc798f8445be2b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_3d3aeb5bb32a489fc273e8ca839343c2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_fd07920044cbf1b6fcf3603ebb166c6c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_c35873a60d1144c83ef89ddfdd31f2cb); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_deb75cc63bfa66b8221455d806440ae2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_7dad59838d93adf969c56acaac791f07); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0008_615ccdcd04fa4562fdc78efaf9a3e52c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_448690263d84211bde0c1cea11a48b11); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_17ec4a11af4369f9cd9a79ffd22faaf2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_05d79f89f42e692717a39a29997da2de); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9c83230a4aa511fa0866251a8cd1e0e5); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_52a0901ae15325b3f8f222a7da4e4491); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_7de64f0cac233a80cd163f2cea7cb779); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_267db5602cbea94318afb38924733c99); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_773881e5943867f63bfc186bfaf4c1c1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_dcd5565feb96fad4810ddae84dbc06ac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9f05e5d4ed89fd9f8879183822aebe38); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_825b00edb1b43df6643231c638a67cfe); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c600a1705c8565072162e90055bd7adb); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_854d6baadd0cebd58435d5cc29185337); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_13fbd5c85aa99a91119179870858fc36); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_d69fa1ea920d7f3f8644d49fc0b6aded); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_d9a28dc8353bb61b24670fc251b55173); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_41d773bbe4db96713017f5c8ef306195); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_d444a162b54c96890064f3460ddfa9e9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_196ddefd25f93c78224271197dc93b05); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_f320700265733d40b644400d2731e9df); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_8a643673abce113f9fea5ce1702250a6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_bf737e0ca2c3b91f381fbcf33b58c0dc); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c79fb6c89b83a91235030370f6e24c8b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9274e2cae2bd02af387c9ce4dd6dfc10); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_427e4e33f79cba44d05a4674943c0473); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_01f5f56161c2ba3dc565775f6be2f8d4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_fea3910637ab9db123d8a25f2431f80b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a7936c92dd92ff855159251fa022a982); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_5b4064ed641abac0e70005efb3a4ced2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_03a7f08f8d44561401428e04248ee99e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_265661e823b3aed3ebfb5ca65c6bba29); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_f321bc679577aebbd61ad4f183fdba3c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_b30963c7d0cef922992c78d9a6f4c315); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_68c6b44c1a5e01b29f32092c5f1d6cde); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c74888bd57b65af1a21ac4d0e5030074); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_cf6f9229ac6b3ddb490a5b585f25091a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_41eaf4b14e8f3982b4305f1f8f903133); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ed1ec44050b5b1aa1aff9aec0885cac5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_b8b951709b027931b2ae953b34b6805a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_1b0976ec19c297da5164b0249b0e47bc); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ff474c2d7415a5c1c988cb0a1b03cd93); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6d67a7a3477ea85d6035b87a3ae74f94); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_4b8fe7a4bef4b14dbe843b5a6815c43b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_03145e837756150af12207ae91cc5699); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_dd667cce68b8575b61aaaa28533db684); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_68e4a2e940a99969f86eff9bf0200ae5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_e0d3b9c4b9a84b90126ca4e10d086879); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_08a849ae5b0aed116fc89ad126ee2ba1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_2ced1cf7ed9e983a5dbdfcbd821607dc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_b6cd19cf9a148786583dbf6fddc86891); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_deff6995f2d3779f620beab780c025cd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_7e1a1c65bd708062952b2fa850109a58); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9d60f26458c49e00bc5fd43d2603e619); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_3e05bf0f8875f261f4b42b7bf0b57666); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_051f43d212f7dad36785e63351b2e15c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_f52251e07be183ace8bcae49b3eb1f60); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a0e97730c95bb2cd3448528f337ead58); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_e42bce9777e980fe284e9532ee788712); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ba4fbb7dafdb5e52846909c46d3978db); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_f67e1bc8ad49300f830c4db205468d8a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_dd4cdd7fc99cf797611e88b4136ac1e2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_05e0f0214217bf6dd30984948b2b6388); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_aba3918be493e52b0f4c7a39f4dca5a3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_504dbaf2db1faa95f1f4b20785cd3e3b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_71e6bf868d5a5f2785addbade2e3168d); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_40e6f60f31bc61ae5b41d15d9bedeb16); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_b46cb5ccfd4cf5e30bef58f218e7a7ab); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9cd76498f7361882486b6d26fbbce31c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a9005e9e8fd249ed58a31bd4f3cefaa2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_01c3fe904d27d62954f4cc41a688e7e3); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_bff64378e85795f8fabccd7129574548); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_8bc59a3820b63b160785ea321b461a7f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_0621ad4aa976a74a406bd35b97479bba); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6fd8c7ff8058ab175550eca5ed66d46c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_0d92a7061fc421b6ab64b7ac074001ce); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a715df8a59c6554c81c3f0a12913b128); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_e2480304f6e52f66531a98cd923184a7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c296c10f37d01b468fa5a6676f0ae7cf); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ea59d342c0977cc9cff54e34cecd80e9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_d672b036c3edb3bc6f356de7eb6fac7f); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_124ea9d194f98ad2d8ed38017c54db32); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_bad231076924f455873af75eb0b600cc); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_50aa24e8da7daf6237997b1479013bf4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6551ce77766b695dfca467061ed69cd3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_42d966ceda8515bd592f89ae3ec5bf9a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_0991fbe5bdd96ef64ba2d9e4e193cc98); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6f1ce97eab2efc4177f7a2e83aa9b03e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6cc18f335b4f16ab71ea937da7bbeefe); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_f43f8c9b658de04e9329743f0314759a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_606569d605eb150571b585d33a7ac035); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a2fb1ec32da13f4447d41d93a1c8f712); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_fc43713f9f34df269e4de565e839a72a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_70d68afe4880001b7820bce2c6562ec3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a4184b99a265c878f5a78b178f8265e4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_833f64043d5880be5449313155b52ad3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_0903cdaf440c013c89e41930fbe6e8e3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_e0cb9b3511ffadb5c06832aab335700d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_fc0e1e5d9ddbdc43483ef6559c753897); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c44f50a0aa4c4ea4213e607c52ae032d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c55ff203396d2245b2978e1088c80597); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_eb39a4e3548820b6901cebed9b2cce6b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_dff11d8a15f2ab6588886efe402c8501); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_8b66ea616b172e9ee49fd7acdf30b938); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c9eac81b72ccb28661bbabb3f035f727); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_1210db1324cadd2baeae2e3872dae2c5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_c1afec2afbd894e5666ed8ff5cccddd2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_bdbf4466875acdb7dd6e329f0d22369a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_86034fd5bf9a19faed074c390256416d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_5a727ba6b819cdb39b9b7c0d12a10805); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_27c3ebf55b07bb44520d8c49cfc58f1c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_41d5d443d3bfb31a75284faff5774eb3); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_cba474b5c075f5f15ff51ac5f566807b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ba8eb78fcc8e48cdd4e64d9b9eeaa0ac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_0e15d4fedf353ca2f5654696254bc8f6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_9036e8ca98f6539a51570632bde39066); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_5ef510e1372bfdb56fe02f981a8624d8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_ce0833b27dce8a567edf9a552aa55495); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_a9c99c789ed3a7ed62eb57e83d7ace3b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_3642ed058ef985fd327d3616060baa99); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0007_6754fca0c4ee8359d1c34dd528128924); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2787399cc3bf9f516640dc94361e2e52); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4dd1be3407f9bbdb73f25370671091c6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5bbeeac10db407d89fbbbfe688174dab); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_7adcff1f6fdce3824e8a593ed868f703); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_36f3872f21fced8308f4e8b4173476ba); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_20a62cde5f5ee06f8316e1918b784015); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6702d13c01fde7514f9cbb8fbf81348c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_cdf3b1387e4f674ad293aaa4273a23f1); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e475241af1c8393981908620758d5e7c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3539a7f57152c3ddf090b34d16228572); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0e4f10ef28366c59fca887d987e448af); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_fd35073521bd29b9032b8e6343dc1582); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0b9f1b7a8cddbcc442bda315e1f5aeac); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_08a845fcbe1fb6d2263f2a31c7fdef37); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ab56b9fb37ba66c11b502a6aeb250a57); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_75a5bc2f988531d233f1385129322def); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b82bb5e09eb92611efc96cd005703cd8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3dd755892540f437fc813bd4f9c6e76c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_91a136fd0b4f47d3a4d4013f8095662c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_1b8c4a926a314136f4f5b3211228eac9); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8613dbc51a7af5956bbd224922fc8c97); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8f10b08e9d138b1d7c3584d1ecc1451e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6f5a532ec7684dbd100cad131c0e6938); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_914f284da40915742abe2648432506ea); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f1fa16fd964817de8b5d3eb8180bd696); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_d2ca153f50787fe8ff5d4d9d15df7387); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_600da2dbe24b60de75ed76e0e20f62d9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_dbb137aa8bd155602797b2cf69146c34); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4139e5345038cc13f9683ec9a592592c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_75faa058b6b2a57ebdaa36b174a9d821); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8dc061cdbbf61be8784e703979185c75); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8c2573bad4170e80835ea383cc59e9f7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_34d5ca3511ee070149519075ad8a1d35); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2470eb657779092305966fce1e8bfcf2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b1e6abfffd742d207c359fb991a8ff2e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8e6b467944e3ea6eb4f2d78bfd8f4f92); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_d079c9555b86aa5a077b492b67851ebc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_72a941a4cf0a009c149b56590a35c24d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c7bf94d6c8d58e7532743914c4970352); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ba4f7230813069263da91ea974ddb921); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5b01e61de74885866a982713198a9e95); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_179505ccc3d31607ea1556bc528493f7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_41fdbbe7cd36a1200433f6da51add330); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6e8d1c43f28b1bc41b6db4315b0c53a8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9028f55e39a92542045cff44d23e5eb4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_af905a4ef6af5f1b78600b6dc5f9e989); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_05e5e3b8142e65d950f76c265dff70a7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_18c3ef0847a75edfaec330b73104621a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_82c86c8ca99b8c772822016f8bcc6994); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5d8d30cda81748d699df4c733cfeed7d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9d9ba41da9501587188327a33b681b5c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_41877d4a96fe78c938c8dffea6a8a8d8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_21c6aaf05ebd245d2a81458a7c024555); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b715d13c93b8e3625ce37a07ca6653d0); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f7af9ce566f7bb1681743a44b577aca3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0aaea350a70215c2fcedb8bc52c75009); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_92655dea42245236225225d1f8fb9ad6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_de66839b38d0d1407fc536262731eeed); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_63623688b7200b99d53ac1d8d149f446); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_57986a7bada864dc9bd27a63f6df85fa); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ecac1f5922c3c6d122112bbcd2ea5417); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_daaf5cac24f21272a745e2120b9fd4f5); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_439a747f420d1b7d8502ab6eefe9a4e7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_792735c1a1541ec6975f0bc91c71e00d); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_47f531e2e2de8e5f660003b96048da4f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_606976a799355cbc2ef135455ae03a2e); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b230083f64b1ffcc0c6bd334f364f4d2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0050743478792316eb42a9c00110d40f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f18377cbcd9fdc9c278c7804f50192b4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9dfccfcd542b2f94cdbcd29ab48c84c2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_039f968629c44dd3bc1cb7ddff0d70f2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_bfa31908831df498a7f6b53f63ad2cea); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_976d373277dceea1dca120722cddeb69); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b546234b9cfccabcc950219c64db62a9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_60cea439eb25c3dd4c5fc08377e33e86); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4e47c679bef09f514d12456151e38722); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_631a09a8b027fb27be1314d88c6c90d3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b1b8b2043027f1c9d648d17f54039e63); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_920646f9f9cd7de1ffc31a5841b4fbfd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3f35b9af9c8ffd4768e33f0c2f1a955d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_44acde8afd404ebe019eacfd337dcec2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_db09fdc7c55f1435e80abe7fc49f0851); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9aca8c5b7def495d806d3365865d62f5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e6577aa3127a1ce75e3f57ec87f110e5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_960aac12246746bee40e9f7def10d535); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_df7dc788652286f950681595791161bd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_17f465977dd7e9292e49f3d3750b78be); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4c7f8dd66b8f56bdc4c7d957a9350361); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9e56d53a758766c7ef8bf18ffaf30cab); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_98d96fb40bc9afcd0640826a472020ba); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_aa58b7d423443aa26503c4665dc3743d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4a0dfa7eb9abbb979297292dd3338e12); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ef9a5be1db1823045e742d093b8e1679); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f0703dc6977f55ca18452372104a81ed); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c3214ef9cc08c38f4c20662dcae72311); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f5575c07a41874f22f8f0b5ab4686728); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c69289d8eb7ab88d6783ed0bae845fd2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_069dfdaed1b51f9cdbe93d4f38bef2bc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_bbd5fd8057adff78016752e1be0a9d10); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_55b10ecfd0dde9d7824b98feb70a4ae0); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6ebbdc3a0d909a8ad477a64e2c788e4f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ca53d6cee021350c5a67cf7963d3461e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_33528021bc747485d09ebf30531c578a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_1a5ff64aa5b26709b02c47092ec58a5d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_16ee6fccc64780a2fbe5db4b03c6d491); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e25c3a362b3347dd39e21d3465ef0049); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8a011e4074c65b06c74e7153bbde2b24); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_a41fec55a003e4dc7e9630066560fbd9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_d6dd182e0deb72b02d08aca8d6474e39); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f476294ed42363d7c2356114aab61ba3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b9dfb6762a491c67ff005532fe4150a4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f2cd072da115c52699051ea5b4b4185d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_d50ee7efdb57dfc2c826d4e82c8b2f6b); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_45746a7b059c5779b0d5f2f07cc17264); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_82de9774803a1efd8ce92001a7d11f55); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_677628441b0b6583431e141a81e425a8); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f412a098fffad4615edc59c07fb7df06); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3eb7339761e80ca76a5e886ea74b5c62); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e18b3805f6755a4369577290185cb406); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_bbc085109b2f0c1e0d02cf7fc8e946b3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_a0bfa772ad0992f8541f658be51b7492); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_560510e89015e29331831d9b1ceee341); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f9ba701653a8c6eb07cf8297d1187aa0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ec4e435a7e759fb5447e7b6bad64a8bd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_467ab682e7824a422bfeb422a169f5a3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f3e8ae2b71ca5abd23475495944b773d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_08c836bfd3ac7defc7bfaa87a18a59ce); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6aae8e835f06e746f9478a6ec320ac1a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_07b442328d533295abda3d01834c7881); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_1d1cc8f0334f9bf5c55f5fea90dbe9a6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0d1a47c831930910f2fdac0d91734bce); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8ca52485803d8efbdfc5decc0c408e15); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8f3fbdc92798a9fea1f7a0449454dd1a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_00546873757c60f5393ec281b0e23e2b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_23e6077552c6edf7b4bf59bddcf4ed1a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_7a5fa0f486ff80bb9c7fb0e446b45447); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2c3fbfeef83ea24c40a3fb6dfa7d28cc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_817a12cec18128ef9c600859b2480cea); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2f219fe5f82e2f01a528964b1c0393e3); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8153005fb418b0155b6e8959aba96479); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_f696e01b39fbfe1dd45edc9adeb4154a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_afb64ded6adc5fbc26c03214a30c979a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6cad2465ee1d7bf5c6248515b4994362); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2e099fd3ca48857fe0b9782e2c4be64d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_7bc596b238c787350281f2424cc819a8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_587fc8c53a8b9a50254f72cc34010eaf); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_480cee4bf05ec6c4b1cf7e4c4f5aa575); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_30c76ccd70969954c0ac7098fa6c4855); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_a78ba5b57311d51ec885c88fcf62a0d3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5802b2f9aac2139509d07fe2b7e4c96b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ec5f3f3ec9a278b63dda59c8c26ba82e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c9efb9ca182419988da01c441ce0020a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5046950504d8a3b58a4241a8d69a199f); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_70a752217140ab71d4273f5e892b9aec); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0269e835f5f63b7464e26f2e5ced842e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8830adf7455b83971bea3333d00d7582); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_fdb81facd138c4a0705bd634d6c5ec20); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_7ff2451bcca1dddd640a52a8a1d0b6e3); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_12a9fe7e3480d6e022d8342adce5e883); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e6441cc632e3a36a18821ae55934537b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8417a86915b2c41699e311e1f9d761ce); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_72f0714b315e98e5b6770a7be1a30036); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_06a0e4d619b983401d1cf6fd57b63b56); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_254aaf2cfda8272e56d4c15afba437a1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4186142f3f5124cf0aa4c9fcf856b2c2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_62f57ede9f4fdca4008bb8b102bff7ba); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_da7b76d23561eb1cb8672382675b31b7); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_71cea317b809e8f60b0684946b8caa5b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_61759473127617b1d4036288e4a8c0e6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9e9e48a167e18445a9e409afe06df946); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b023136fe1c19ee81c3975c762b0b8e8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6de7268e0b98d0ef7f3b4240d2068596); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c5ce6bf14ca96cebd40fbe21bb29b7ae); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8c4a4d51c5e2a5272c519b0c73f39980); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_da4c4979dc91a5536b453ae9de6ce031); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e7b80b1f90f2a6e8ba69851be6bef3e4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_346ac8866660b5adcd80229f1085a8e4); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_4ce176a6dbcbf31217f4f29b6694c68d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_99962fff9dfac7f527aec05bd41d626b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_22ccdeeae81eec7d535dd685dbf5948f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_6bb429c97c0ecdbf2ad34c4ca1332c21); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_e771d868b1074b753fed0d5c9a3f9c23); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_59a3746398209af8d5abd1b9cc6e3ec8); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_728989255f27508682afd8f3bca8740e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_37e713c9732992f1fb1e77c0d97b5eaf); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_c36a0972036bee0218396f521005e737); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3cd50bfb5d15730a5bde269d70e1bc1f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_904b3445479d8e81ff729ad84d948d4f); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9a86c10ce63fdd3a02cb036abe9bf7f4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_2c7c487b92de5c26457d4c323d171865); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9bbef3126b2a214738ab649744f1c519); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0de899839a124cfec12eb9923802d50d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0ffce4a9e5d1475e6a58d770dc12865f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_aa8efeca905577f57d9c7b5b271a47c4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_66fcfe34b767714399178ce8b69ff974); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_b1193c77c670778f320d60fdff28003c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_3c68d1110b4c053065a7d0b0d6bb71e5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9389c63162950678262bd1fce5c0b4cc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_9c8b221b55cef0ba458722e6774a210d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_cbff1b1901fa4fcb5439f11ec0e5ab78); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_02ecffb5ba3ae606783c004a3b6d9580); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8f4da85d1f4848b9967393bbbfcdbb5c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_fc5da0946d81417cdde6ff2b10d5ec26); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_8f1d7cbe31c9ca537a8fb61a5f9f935f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_5882217b398b67e9b62b539def98ebc1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_0021dd584f3a247e9c2e6708f1a62eac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_47d7989aa77c482fa16a6787cc24626e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_ed3c586b0f2e0e36f908b9487bf6ec5d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_7978a2c2f434e7dee7bf0f686b9a14da); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0006_beb0eb690c4037d495404a64cb1ff43b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2ebd13040db833ffd88ea5a8d2b377c7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_cf0e2362c7b8c8b7b50e774f3e45a654); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_dab320a34a1bcc5b257599806c5751b1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5061a7994654594f8524aa2ca1c1d1cb); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_eca72774318fc26016719a87b0035761); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_92148f5863b6f107ada35b49ff1c75be); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8c8d9fd00cfea2ecf9be23c55594571d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4c04db0f86bdd7695eb22989598991ae); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_497c507b429fdc6223283c0a0a534ea5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_73f928359b60acd4f6a4d9ea2d69a0c4); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_777138835a2952632e61018f1c9e7854); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3c4f23fa5c7397da192518024423d3ef); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1df2cd699f63d177e5fb8d9548e2c0d6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bf70f73aa456d2cbf25751eb606ae32e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a1fb2bb0c3e0b4e6e7a135ae626b0a33); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b047b9570225719cb7f183bb00fa4daa); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8cf3e7e6e07bad454a5b234c393a1d02); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_976d603d13330a49d799741fc8645f2c); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_6615ddf788dfa153e093b04b5c0304cf); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5724728f6660a1e555fbcefe040a1726); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b33d9f5f4ac39342663433fb14ff7314); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_65e1c77325da4d0c054e69bd4551c599); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_eb9908b466f20e92c27e6f5d7aa1bacb); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_951eb871b1cb468cf6c766c8c83e47e3); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_217bf44dbf7497d3f13ac203754a7e88); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_e3b29383c07170cb65cc52745a37ddc2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c80cabe511a432e4d8e67766052875e6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3da95f6ecd0aae49adefcccf5ee28c1c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fe91be5d8911c57644753f925082c60b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_513123ae151878de0171e3e4e1dc4c34); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_dab2c71afcbd6af7c7e4e507d436ebb3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0e4f2ce6ee8c6aaba65ea53cad0a7eda); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d6b0f2f4ee5de04ac98eb3f61d53e1be); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fec8999873ae214a7772c4e6cc77a241); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d90d9f2acc455272bbb9d9676420675c); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c05b0bf06ed751b21d3d661715644274); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a1810061d6b6bd714973e83c4b126522); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_71ab8e9bc980bf0d31da9ff9c7c51ebf); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8e89c1bd95e4e9731efc7ac1abba3639); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7d6efd8ddf54ee184c44b622a1ee868a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d6f7537a22515e5e9ca4f620125943e8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b09da8d8fd6084febcc163e9b40d8cdd); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9126fe4aba9b3c897e19df3abbde4c2e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_48bc1b8fb59601a4f7543709bef7460b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3b53bb12eb2356991a95238fc2402972); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_95fa89461870d20a1aa400115550e236); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_370f70e1f7779d2b18c16b77796e6ad4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4a019bc6462dc411464060a64585b9f6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a5b461a6b03f2288d0023878bfc87868); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_61a248bb5bae07d449dffffcab033bd3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9e017ed93995c8a79a4cb81af48215bf); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5e65f4368961d057065044067e0fdfa0); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4d91aa8a5a995c9a419786a35a9f7da8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_091322757a98b7a8c1de57f054e6eac4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0b7241fbbde4c659f394f9d5233bfee9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f65dca199c8c7e749b3b1753cbfba37b); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f8c02b123b7eb2a971596b935c942d7f); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2e8be41ad679566886fbb5aa5471154b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_650d3e7f112f289ca18b3e91964181c8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b1c239773538f5ab7b4ce96067f7a245); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_910e13cc5d980d54df51383c6deaaf0d); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c4039cd70e55947c4038bbc29a3711d7); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1b7dc36e8822afae1c8bb433c828d4e8); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0e1cf0431fb2948288bf31c84c442195); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_db318a810481bb79ec5c2f495823168c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2a7c43b8b90dba2d44286b98a9b73df6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_e7dd37108246f62508c0e73239cf6e7f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_81570b25b723de4082fb023c40e4ccd9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_436b3b6f6e649ed483a53237c3775c95); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f2f5ce16d0204a3003a702dae993a40f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d961dc8e20e0cf7d2f38aed674c88d64); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_633fbc849b41ff3a78e8bb75ce2db424); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d37e6af15248fc04337ce571fb477662); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bd2324d2dd9fc896512c6a5359795747); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3c446fe7d7a33ac81540f3f0173a2617); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8133b2ae04117f0b72b3d8ccbcc77744); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7958686f0a2597a7cd7ec5775c25b60a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_6a9a879c8cfd1e0a91ce0222094727c3); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b442530703cde794bb889dcfd2337519); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a47ece9818d8ba7d83c3dfee9cb5a5db); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_34bbe33cab14fe352dc7f15d5b8bd6c2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ff3f6f012645445f74f44e8d5a4f7688); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8ea2641f961a5c26a90da3a80a453961); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bb1dbd16911277023686c10027f764a6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_cded940add53ba9d752aac5b749677b5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_afdced595d6aac26d33b2396172517d3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a7770d458fa6e73e3be636bc9da5a38e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d2d66cb76facbf1d1b7a9c37bcf940c8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2b35a5af01307fa72cbc83a6a0a89f3c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_781918f6c2f2d7aae6eb6fb4b98b8ddd); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1c6e0df7e9dde9612528f8d32ba1b892); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_26c5be906ea66f39030cc40db09bbbef); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_76190d9e70ad091b7248c8287842d7ac); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_63ba7a0c67d74cdff3b027a57f010878); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_39a8a13222fe1835ac5f0063826159aa); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7081086f07db237298e2914bf8f98d08); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_55376846691a6cf458c674ca9055c9e8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f0c2880fecb021a3407210e94b5da7ce); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_dccec1fa0b1ac6076c36fab915923e7c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_20631a232115a4d613bbdd90bd786e91); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_db9df09b94b4eb586dda5bfa0be9970c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_76bf329e370e4f9a3dd99d83026a5478); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1fa6f2b326e4409d7c306daf844e421e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7a8098e8711012274339764352391092); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_95dabc404203cf349885a9d084b8cd6d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3dde8f11372ab3df00c06e9cc9f532a0); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0c113831bfa3a7d8a21e204222a252c1); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_57e63123b079ce33af8bdadd71267625); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_43e95a1c134b51155bf9f0443e1d2244); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_e4c37a5e959f41b3c8f64b795b14239d); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2b060163070164b37e556fa965a75b22); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fbf26a0e2f1e9287b0f1b8a7673bb60b); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2bfc2a716810b8b8e428184e9b039b54); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7f384c69fc510e294062b2503fa418b8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_15a9aa0a2fece20c5d53b7084ca5a21d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2d29a325eb959e4d792f3c431b67bd2a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b23aaed3817fbc91a8f30fb2d067d608); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c86504accefd86b5cd487e87c469c6c2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f9d88387943c5f0ea9186ebb77de41cc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_589a3cf2c8151915d84b023f02fc752d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7536f936cf5efb906be9d5c9c4b9b797); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d2d1a3dcf8158900936c7f1285b72256); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_43310fbccf884406bda0e2d6982345a5); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9ad0458e4f0ef59ff7fa59f4c0837b04); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_95bd9777d6a5351f1331e1658e91c2c1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ca248c90d42e70e795469281c93508c6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_399ec08db4e9d153838525614edae502); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_646f19e9526b1b663e651ef0e2bef0d1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_7d9a114c2369646553f1a24e06faf975); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b486278e1b0983217ffb884258fdaa8a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_8b0d3838be5654907168a76bc7af037b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_be4447c4210552d0e2d709e04422e46c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1c904363cb592579ae624ad523c9e4f5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c19f9e19c46b1102ea0299fc9a1b38a4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4fdf31a13c8ef55acf9bd6871ef8b7ea); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0f04d6a673a007221d6a5b614f48b38d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_e5655e77364e3f931a3ea6875c01a0e7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_24a1c25e2a27df99341328046a46e466); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_95c51c7c23affedd2c9b94e0b9376ad3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0ddcc5781c2afc1f94dc5b82a654a243); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_455bce9dcc99aa7c0fe96735f9c1c93c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_76107da66f5dd1fc99f9f8706bdcbc77); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_13d3220ae21dcb3631a0e0d628adb44b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_21c6c51021d4efd37fe596f5269fc7c7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f1dedecc1b8b0f34a2028626b0c5f9ed); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2811e4468999fdb14616de0d05b9693d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_61de1cea77f9a0a6f3a05cde41d82c0a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bb0588c994e9ade1989ee58ccd0341e0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_585bb225ed4b545c00cd962f5127ec53); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_50e95880bd580ebd349eb2d5d91ce2a1); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_15b43c6aa30caf41244412fdba5c982c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3221787875e443cec88a71604e26cbab); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_aa238e8cc38c595d62af9027950c4b6d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ac2671a78c195e70a534feb05d09ed8e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2bbc56277fe4dd7194f1c426ceefcd73); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_036cfac2479cb17c204aeee52931f723); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5b147a73931c3e3f8e956b0ee94f9d2b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a4f76c23639740e140aa2cdd5a506a3c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fa8c6348b3aa58004e4cf20cfecc660f); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_cc1e37d25022b50f8aafb1efd4be0001); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c670b45682d64d3a2f572ae08687b1c3); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5698322bcb20d1da317fbe9951ea1ea8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_be31557dfd123659214c16d5f27af1e3); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c829a7b8ed999b05fe10801538e8fc57); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5be0a90c30935b8c105dc54a4cd6009e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b971b868aaf67fc950f128ca0831b10c); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_67ce16f2bfd9eae5a2cb635dbc252281); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ebae4eedb7152c7969b661e0ef346b25); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_45e8e15a42cc8f475da57bb55786ce97); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c159d592078beb8c2d756664e81eae2c); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0a036aad328e938d09bf4e0bf0b1226a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_731476203a8004766e146f3ef7ded738); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_994f8417c053f4d92d2afd228765962c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_11edbe1067aeaf4c860c11793036615a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4309264d03baaeec83243ddba351287b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f86fa22f0a82b786389b11aa022dbc7a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ff3bb63d2da59b67dce62fc09a7068e4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_271d69a36bda50544d2c8d10833b45a2); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2a7c145d4ad450238d966cd1b81cda91); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_efd9a40fd7a5179c2d82d3a2361283ae); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5f52d178514fd4128799774c187f794f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d10ae8f63596ae79b0e9b3e15c62b219); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d8e90d85916412acfb7ae06a7cf3cd72); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2a1026d851dd597c88f9d8d7bcd46f08); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a613f56e071157bae4600ae04c11f63f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_726062f72cf8c22afcfef969d65edcd8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_69e39a5cf47c21d10de2940a468498a4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c3e01796bb9f89d16ec7919b28547cff); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d7c49a9c69464d5593c8d37de766beb0); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2d5132626aee07585747967a992314ee); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_b9f4121827f5ff1843626bd5f67092d7); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9445afc0bb6616b0ffb022c3955ca01a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ae364fa0447a19bedbfb479ca3ae419a); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d18f395efe4408da08427fb8cf6b0f8e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_cd26b941286fbb8d42729f5e019f4e18); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a1bb57afc384646145f6f8acfff07ab6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_34c6fd3f84dd7e5b474f79f57e7410d2); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_858ad5997c7a1ecd0c4e016da0b78cb3); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_451122d6894b9c1a4dc9870cc5c9f729); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5d88e7ee113b13cb24d3606d5b95737d); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_343af99016771b4bfeff7945efa2b0e0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_3655878bb6cdfb4990bf98ac230ba6d6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9bbe32024047c4fedb1cec5108eaa882); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_612cbcd42b1ca988c32b19b0ee60fe8e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fd14a18f00f97dfb88beee48134b1ed7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f305594ed0d6b427fb2fe6e450424460); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c391b2ea017883e989dd4c91421a2518); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_be38afe7bdce8819fe9f4816cfb05e2d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_fa236057bccd3e2547a2fecc91877c31); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4b313c6d9a99a43b40d16ab65bd356a2); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_30e183bc3b73be79a53a3c5ccce7ceee); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_9820f7ab194cfc8681fa530ec810fbee); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_86594d133db62770dc866e4b90a41484); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_0ce7c7f80ec4cb8b2af6490440fb797b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_6f9178fa3ebe1aed2696357cd2afb3a9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1a8bfd6bc3b5d039e0d1ce8df766903c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_5e64d450c4d7d27074726a3c99ca4467); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a03e3132f8a1d5217b2ffa017ef0e471); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_a6894db2a6de9012f66db439c1729127); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bf3dcedb268236b423e720a59c72d244); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f4bb60d74c1105c14eef7ba20497175c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_59952860d6c23bb78489fc55e73bc893); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_eb6a678f51255289dc34981b27eba51c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_063896be35045b544dfdc8bd6e379acb); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_522573917a875da9c8a8afd3f5bc8e45); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_86f6a80584134eec1e22e2bdf0b8b928); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d1809e40a550459024766fd7df6b027e); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c19e3acd7931b2ef3bd181d17c473d6c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_020e8fd0be2f1cfda622cd7fae226094); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_6134bd04f1e6222f7503813d44eae9a4); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_e1fb7970bf6eaaa98a1adf239f3ce9f0); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_4f703912b1d28a906d4289a579b5e53a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_96e2182da02d44001714b6f444e74479); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_af6f1f46ce4181887d43e39ba0eac2a5); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_58c6e55019354e20f94dc97b5f01710e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_11ca78deb33bdf89e23765a32c4587c7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_20d8d229b2484368b368c970e6deaaec); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_172f2782b0ad97eecc7f56f17323847c); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_70c485a13a46b9c5b06e852ec8bb7c9a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_366a50e5d8552169108ca9bb15486b99); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_31d435075b7bc220179967c4a15393b4); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_f7594aa16f0a87f7e38788eb9c5f8e4f); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_2b9c34cec984ef83aa9d7eb139b39e1e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d98ff0a96153b4e2dd5d28a88b553717); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_1f8c595ee461d1900fce4cab06c9d1a3); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_d23443a383816cc688a8e378828508b9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_6d4a973194a68c94feb8c4d04ea3b1e8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_ba9fd7d66463d2daa1066b70ade39d1a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_bdcb0b9943d4ec26d0fd0d86f47ba316); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_af83d71bbe5f5d8c2f119446c9ee2404); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_c31c7d4b52033249a8a1d888dbd92671); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0005_91747f9e28ec013939bb19940a53742e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ec4d0075fd0ddc8b390599c00aa63ac1); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a602fc8de20f9667dca43dab54118301); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_25d1043092b1d4fde69e802678864cbe); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4abb5ab2b97be2656b1c86356a434af0); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c4837d1be7ab7d53ba06b807636c0d82); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9e0ff90f7f6de1688dc2a4958037811a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5497ec2107ee491c121c6b9957579ab8); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7c2de3c593da1d0c0ca997bfb7e1d325); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1b5b9005d979312b723d3a27c0c65db6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c130e5147dbf82d5c4a8d512eff1636c); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_87e9112616a76f4e1f0250177acda175); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_09adbe5008d26d1ca9b6390889d3e56d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_393cb9100ac7954d3e25afc8842b5dcb); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_bd27c65e18a44a6a66f0f10f280dd62a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_457719d02ae8d6b7caabc81e5ac31d57); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6afb89b2559180fa651920ad74abbb9c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9efc0e3a9cf918d9d10523202afb880a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_dcb06d9a6a9344136132f0325f437e0a); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_421d78dfe11851b6277447d6fc3d762b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_62d2456de8e572dc055cc37196ab7b20); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_3b5220a9e11f03a530f5405a48a3583d); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b214fbc2704c5d1b1313515138517b76); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e096a7508a321da18cfa19aba7f60b47); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_da3fc41871569e3343296bc91c49ae46); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_62cc580731a221876d4d44f4e853a5e8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_54f6c56c9f04d106281a9ec1fa12677b); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b67343921e6e00d89b71c1d31f11f704); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_2d60dbf9baa0aef32c08fa102e3e281f); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_420b5a970fc5a50894e8f4bb8ee9ad8f); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e3f29053589874b44fb4d5ed1443d8e6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a5ce8739821cabdfe852f4ee9906463d); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6dcac6cce66b4cf60283b63f67f7ae42); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c87c98d8350c1939120b60c35c42174f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7f966bf081d9facbed0cb718606514a8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_42ee2395300cfbfe759565dd572359a5); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7c9660a8ef7ac5680c6a0ec70aa74f7e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e3fef4b4dd1fbbec0e7621ad79e33759); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_41a10661ce28acc8f733fc3a3f82f02a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_721d0a3d0378f788e179d749fd0700e5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_05be8063e1e61e6a640738667e26ff1e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4a288634e9c34e5ea3325248b34c84d7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9f44333dbe686be726912364931d38f9); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_002ba8ae7ab8583e9fe9819c79f7f6e6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_66435cfe8b3a54251d7169a90c4a43bd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_575b7d6403f1859046727c7cc7fbb6ea); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_438ce4e1141c53dba704cd795cc91d37); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ca46d67dab84b4a7356a56ad2734fb85); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_97fa4a641ec65c1c85b776a1ebeb05ea); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_3130001e0037ae9b22700b7a0c290575); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c971b4273da1e06a47abc0e7f1d294db); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c9b76d5a245e1826763fdc264634a414); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1b996e5d7c48c35626c06ddf2a0daee6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_684b88bf9d678a38c20afbfe31f4589a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_cf59bff90cf956a08134977785d32526); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_bef873ac3553cf270cb99c91f8ad61c9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_f9b0544f2b53b646873f2d45ad2a60e8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_d72ac6dd8bc2e2fd59bef9bc1eeaf860); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_00169cf72ee7ba4bf707f0ee573f56b1); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_31d4428785802ae57a327eebeda993d1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e0eb441d5b5198369a3d436dbdd96252); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7b829e6cfa007bc3c4f5d8ee731ae57b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8497088c8479c1181a75be21c357c192); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8e129f1c2290819ce637f1864c5db12f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1492bd1459024954efa63e5096a5ee23); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4d977f1346507e064ab5be6f55e3c261); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_0b8b629d5421fa2672c401e0c3a208da); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_86d48c27c6f553db31d7023a0618aeab); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1f7dbf76ba241759ea22fd9e67b0c97c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4dac992fd9f94a93dffaf95424246fa6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6a4c903cebec66621da83b719ab797e3); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8c1f1ab2627a320478eb2bd9c49efd73); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_f679a21ffd1bb700d447688e1f22b988); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_3e0d6a60d0b94f9a38bd6278a9aca508); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1d4652a872051fa3240ca534bef0a43e); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_994bda6d7619bde46104d2eab7ff768c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_63ada1d75739c374aa86100675e5d5d1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_052895537635d2f5a24651bcc48ec975); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ecaf12077437db7ea718765cbeceaa78); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8fb2e6de52e01bddf84ba9aa01f20489); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c6e51aefd95b2286637499cab5834392); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6e5b2720e60eff99a0177f056c13773b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e480d4975fde5e93e67c0901d131c42e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7bced5185e944b855170b3b92901575a); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7e1db962b81a09f792909878ea32a5fd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_7132619a8c83606d629fab216a12c395); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4272cca8760969ed2322624711f94407); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8098fe71a9a9219d58e590ca89192425); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5d17fc1b7064520ae6022fd2f62c2bfb); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6fcd36e5fdd1db57ef48d5bdd4beebd9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ff336118365ff6fcd2a2e761486717c7); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_62deeb4f7a0c751fdee63cefb834691f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_855ad8d8a1e77072d4e79e775df962d8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_59eff3dd52fe69e24ff9a0e5fbfff105); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ee03ba98204ba8f2d9fbb2e14d6a728c); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8fef826c40cffbdc6a69c70b5676c824); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_62083746252ebfd031b91ac334e101f1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_2dfdae3f68bd27ada2cc3a92559282e7); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c68aa300799a706219ed2f7ad15362c5); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c948501ec7e870a27ec7845843aeaabb); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_68b9285a81940136122cef1c146a64b8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_0b99abaa61e35521eb15190ed9f7df0e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_81ab8fa6399bc4cb3f8b4ceef191fbe1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c44f3469765232bf2f3748ac3491ee68); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ab55d13e991c298837fd0697904e660e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_573d41654139e877dea050658666ffa7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b560d15b65f988c4960fc412191c25df); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_72837a807cd7df18bbd7e22ee697c816); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_abcb14740cbc0d1a201f9176fb24d825); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9adf723bcac874299a491196558e7dd7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5f6080df6bedd5e86600c0719613d97a); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4f64f93a4e925dbd2a7598b37bb8f332); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_60f92c1f81ec8480346add9207661f84); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_f33d7e93149b996e3b7d599a621eb41b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8480b7a2f4b03ada42be3b2e1e8ee03b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_2a24acd339cf38eb1b877ff6e36477c2); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b702ceb57e9a027bb97820e611754bef); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_be26b6411cff1ae6721379214d087f9e); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_daf8acd200daae2412c1313c41d2df26); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_089675b7ca4c9fa833b6515226a0b60c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b454cff4daf9ab30e3870da405d5fc28); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_86f21eab726445da5b1ba4da19a01e24); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_b0e5c7fb09f6acbdfcf6e0d4258d4a9b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ad424cf88ec9bc306be6b23b4034be4b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ff5d54172f489033d905c69de9288f35); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_86340b112896f6992ac55fb1bf8e74b0); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_951e06b10a39de1cad7dccca3392ea82); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_74d560cdbdf369e005ce1b8ce93a13a0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a6112ca861997f87f6d0e56ce5b1c43f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a90a33e85b58caec9ced942971c780f5); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1bdf28e924b92fb4795814702aa75ef8); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9050482dfdb75ed3bda0f45ad2cc37e1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_65ab36fec2bb30a1596eb25d39a63d3c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_15e77a0e9aeb76ab2f04aa748ac4719e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c8f1b5b5a8743b962a62682d4a6e9262); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8cc60054c0093e34e396607e0b128191); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_0fd04b4e30706876c57056f2a95c5b01); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1915c1c223496f64670367e8aa3e3eed); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_554f29b91fc8ef3f3adbf7b1224f0a44); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5423a435d10b320ffac1aac27fe43c6f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ac5d94b9c5d9718d65ebc6e6644466aa); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_d3d3167450d7d79e6e587b43b5014c17); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4073315e2d95994c6f3228bd979f95d6); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_270c1587b679b20fe7add8f451fca302); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_51b8c82cc36cae2aa27ea9f5828ed420); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a02b8d7f6b52abd4bf35c4c2828fa170); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_2e864f41a5216ee217adff2fac1e6836); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_211f8a32f588ff83dc85e793d7ea741a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4adfd3cf6b9fb282abf8009ccdfc1dce); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_40b177ca84280f73c468ff30cdd2e97a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e03d707ef6c5cc14794c82dc511c6c91); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1d2ef1d66fc107f4d6470ccf34b8b037); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_3fc6609110f1567ea6fb2ccc9c6941ca); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a93de0616ed3079543f82d35645aaff5); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_aec80ba98ebb4fb6dfb693519c486fd6); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_47a3327409e70c795fa50aa402f552f8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_020f136197f532507b2b739e893298bd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_fbe9055480e6ebbe349358be151b7b87); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_2e6272eac15fc86b3a447774b973a7db); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_fed4efdfca6c7311372ae1d21591bc51); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_eeade78a4548ef600a551ccc3677a14a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_8b02254c6cfcd13726ae202cdb121b1f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_f018b29703e3acb35c78e692c0b4e037); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_537d8d4f3430c6871d023b75293bca4d); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c6260f7de558b080743c3c6705aa04c6); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1666fda37418e8810185898fd16573d2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a0bae6dc5e6fadde4228e11bf5e49dd1); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_bf56dd41a04eebaa9037e34d998af495); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4501f160909b47a29559b718409520ea); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4eafb813d5396d818b628d6737e80d42); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4183310106ce1532189899565b015546); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_e741bd5513728d7fdcc7ea4757864d07); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_bf95b1c0c61b8c9f9cb99f3e15ce6042); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_28c597b75ba228e5d6902ddc68ad7811); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5de483c883d2daf2738219ea08483e30); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4ff0f4092dd79c548ab0e495e6cddc3b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c4754bae9bc3da992a48068120977d41); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_466daa400cd4119b35117e42b62f7a5d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_819ead64aeb08e11d926e74832dc7856); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_4ef020c4936aea427d194c2b96378df5); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_12ebe2850e61f4a8e07fae374063118d); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_30facce362b8d7c83804cc9435f8bf2d); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_a0f6a4ff15d8d9e4f751f67c2efcd8f6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c818cf1c81182b6fe0e9f65adf919fcc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_94eed8c577d89cc2d2f123f92b3e0a4a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_ac82cf73c860771ac75efc0b39b24537); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_d3598002e745fefcab218230adaf46d1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_86042fa2173f60c1aaa263c2c4f6b7ee); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_28f0c23a935aa8023225aad3b1586102); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_418e5123313e62fb4031520b233885ec); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_47c38651973c418d8726100d2bdf9d3a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_789e9da15950f1b6a8a6c4eb8902709b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_5f52c927d92988b4856ae8f434c23f6e); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c921c04aa70e11394e15484e40785309); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_1df93e659bfdd7643b71682374b5ee3c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c7a900ae364616340bb424043839d4b2); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_645074a289d63811b9e9a8978d5ba260); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_249f8105c1428584ff3103b40440952e); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_555cdb27b272a3f2c0eb9f7b9e0eb5ff); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_d1052a967ff8db9233909609d294e184); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_c5f13f7f8f3b7d0270ae28c5a1d96ba9); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_51825a78fc953cc98d20358550621612); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_f39e7f23bc4a7c1bd04d47bcdf9f2b75); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_95c3412bca03d21c7267d3279e2f3312); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_42ed5cc338556317302dcabae95aecba); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_66f15c6b6bd4dc2a63272af2d58fc85d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_9fc12f8e3f7d00cb13a0a0b563ab71f7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_af5eb26d3936b8bf476e1ea60b3edaac); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_84343f2e4ec9d07e1a9dbe14088de505); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_6d1b579154234b2d3345d0013c542993); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0004_06ee4205d01722f7549f02e0b7087b0b); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_55b0c9c15fb89c88e545e8b9afe973ee); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_df39e5dbce330518630020ccd9345df3); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ee7a3b8aa7f629ea0410d3cb5df23364); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_17f4ce77ecc13085127e932ba6a8fb65); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_f168d215225c7302da641eb4505c6093); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_a5e837eca1906a3008beb989e61b545e); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_3ba01671ca6644b664bff3737b516c82); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_80b435a335887493ba9656f5a0e76a14); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1bf343887515a244e3d82298438693bd); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_0040353f748f8b98cb8da172c6b336c1); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_714551a2b9e17d67c6952821e0182bf6); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e00c84ffbe65ef6bc98ba49b16f282d5); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_086984675e3daca69cb4d93a08b03c05); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ecd209322f30837e1dde5cc3c5903232); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_f589cedd98232e73dba33f2b1ee57f05); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1c8897aa6698e41ba80853991cc24839); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_7908fe1d05a5c3552e8ea5c1ee7d48e7); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_895867613f93da07976594c2ea883249); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_739075d41fcec00bc42fc215117d32ed); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e955cc61481c7908edd1588edb5952a7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_6168ba655398488893634accebfceca9); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_975ad90f36f2a8d5d14cbc769a6ce607); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.QueryerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_46d10e7d36cfa34b2e97b70bd8f4e8d8); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_42faf5a4eb464290a69fa9800a551bd4); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ca97cadb92a7dfcaebf19bed8659a2d0); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_fd6abc6efea6a6cda4203b6ca3a23574); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b0d6367c8ae272dd20ac1a3fc8269d5e); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_f71834a64db6e15385899741cd4f16f7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_d2dda163c53609ef7d4070cf89140c06); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4a1ee3cec4754da65ab631f8d55709d8); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_86f0c8258522256415a6791126d2786c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_89bc20456d7b564e816ea43252326002); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_065f30b751342fccddaa383211967ee5); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_d9070b80f72659573a9a1e1ea8b7fe38); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4afcf6b296cf46480b7f20345361eec9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4a6b7427618f96e139a937955fa945a1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4e478d169a0c925294ee737c49c1ac6a); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_58164b71668ef615a11f9d6dd3ff1e07); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_66ae780e5764601ed26339f0ef146ab9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_3e97f3c285abb15a2fe7d2b9982ba7cf); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_c02e55bb6f1875e5ca0c1667579be31a); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e1b111d2c8c09cfe4d5b9d367c734191); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_8ac86baaa7af7d41aebe48b4e0583658); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_576714a6f8eb9ab037040fe7ef742b28); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_cdbed2933e5c0be5f6806d38207387b1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_21f118405ed242fbb797757d53cf8761); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4aab8206dcc15dbead696d4dd7e33cea); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_5a07d33df77df1576ce8067b8997d413); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_348becb0e4f2e9da375f40a1e59807d1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_429af3d54eff18d5961d7d87b53413b7); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_9b9b5f504f5a2fd077ca61e40c8f53f1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_f3e446fd938f183e46a975476439f43d); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ff274753683e1d0adf6ecc55b9be8631); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_67a6f30839ce48a2e9f3488bf5e2053d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4b5d9a137b792a784a0b494e06773edb); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b493216de6d7d70b1fa112fe946c8bc0); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_2c310269b09748f7be39d17b8528f2df); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_d033831362fd7cd8abd09076a328bf5f); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_57e946d049a38df9bcfe6ef5d8598cfc); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1ce67851a70243fd021fc5ea522ae189); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ddd509d4dc401e51bac497647969675e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ce39b6405b9206d7ba9fe332b19cf89f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_66e6b00152e5329959d8d66776407655); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_95d38f7324a5d842e83132f8f4476f2c); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ea65c7fb86fc48081e1d867a9cfc1870); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_3371fc70ba5d416a901c5118687e778e); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_aebab2fa61e210e3c96cb924b4356379); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b2d88477280fe09a57ad676180ae7fca); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e1f4671a79f9b4bc9f39c5b15d6aeef8); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_75415c9d129d1dd60031d39ea83f5316); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_344b3b12f6cdf795bbc480ca3c22d0eb); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_907ef3b14632933f40c5231177e7df3a); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_2295005ac105ed3d2e96ca764c9284e8); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_55276ca2338c0cd130c39fb53ac4caaa); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ddca4fe2167f6deb65f9ee0d4d46b23b); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_24b8a6effd27e68d07d33d0e6d4e7124); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_01d353c46df9c41bbdd5f105418d9b12); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_26a94233bac30452998cfcff3d4ab4e6); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_5701107029804f6e326b6191c81f660d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e8d66097d62d92ca5e8a01d0627ef47b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_d850e679e03edc84e19f7ec96957d9bc); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_cd50db51f3efd4e06b314497df4e8bd1); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_14a98e8d04e2b4b3493f4673295b2dcb); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_c09291a9dc155f8d0c0f1d1961c478c7); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_87c58cadd923ee575179af1ec5f175f5); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_251125afd2a9bfc7c7bcd717ad73c7e4); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_3046c3f6840014e9d6851479421aa2db); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_0a04ef57a9627c79f86b26469450314a); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_214d339ea349168fbe3e0955424c5eef); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_84c6caa443f9cbb04c70721344581649); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1a34bf796737439096bfa6f10d5696c2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_35e2a252407d971a4249ddd0ab337396); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_47814f9637c1fd14af8c27076a6fc4be); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_a35f893b3d8e8eb80b18a381c5ce721a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_485292d166cf8930ff7f1c6d74199750); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_e575d519c87b84d55ffacfb52b72183c); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_613592023fc524d8f895420661ca88e0); ok {
		return struct {
			driver.Conn
			driver.QueryerContext
			driver.SessionResetter
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_ecd9017910d8eaafad4c84916460cde2); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Validator
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4b7e2498575c5cb9db1b7ab70c55ed6e); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.ExecerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_dc6b33036e44bc5045ed542d781b64c6); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_4b20cfcfb6db2db4449ec9444be99683); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.ExecerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_763309d9ccc50f11156cd9497d38589a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_054b90a03b02ee83d5b32a2618de1414); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b1e4b1e4747a4180955083debd843744); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_0980f76b2d444d923da26fc4a1b82818); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.ExecerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_a3a706415d197325957442aa15129e61); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1e0e8df2036eb4f1f761196bb7b1da89); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b9c9b9171bfa8c16742e5f91689b5d88); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Execer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_1b7e307d8021743b2d0cec1048d55722); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
			driver.Queryer
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_b49194b335bfb08661e6e40529986bad); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_35f2471903ecb7c1e4cde7fd5cfe5e1f); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_68f463fe3d9a985f2c57f5f872c7c333); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_d40c325b5125e64636b9b8ea3f1b8975); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_9d0152a9df319a9e7e4a0bdc7f7fa2d4); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_c18fe73278d7622bc6a1b7ff02b795e7); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_8654b308d9ff695cc51971a84d993c51); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_67e361827a469e5a7fd41b0778aa1958); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_23392d2c642abd59125f77ecca6f7b8f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_734a0c923476b7212c3571646d8223d0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
			driver.Pinger
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0003_580a0ede16c719a4c714bf892fdb4061); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
			driver.QueryerContext
		}{c, c, c, c}
	}

	if _, ok := dc.(wrapConn0002_15a5fa7e0b58b2775996354e26f737f0); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_2f1b86788f06ab353f2406475549da3a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_549f199f5c8eff6a17a1207eaca388f4); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_863bdd334cd7bcaefd9d5085e2626cd9); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_d889ee74c9f7be4ca1b3d5f12965aecb); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_36d783b278273e048cbac2a9db4fb974); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Queryer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_b7ff310fb7e22a7ccbb847bf26429dd9); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ConnPrepareContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_0921fb5e4c916e10370d5e0398435d4e); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Pinger
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_21490f92aea80b01cad7a6a9b496de9d); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_e5960e40a7bbe1b5a19c38f71e0e6861); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_1edfe563a97fbd21d92eeed3f2776fd1); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_9998c1796991ed1018aa6832ef10c20d); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_608b01490fc5eeb54196421e44a1d294); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_fb43968160809b260946ae6f55b7fe56); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_20cb62afffd12ef5a3dd67710522c25c); ok {
		return struct {
			driver.Conn
			driver.SessionResetter
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_b0936442b3763765deb78c0f8945bb0a); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Pinger
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_0ab84d337a68f7f945b7dd51c3322a35); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.QueryerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_71ee22b4aad207f24e83c993b46f34e8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Execer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_57475ba369865550e7711feefe48c5d4); ok {
		return struct {
			driver.Conn
			driver.QueryerContext
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_270a739597d105c623663cb03227f1f9); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Pinger
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_592daa1271d02469a5a34fa382f5c410); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_01ce0f63e6a6b920df90cbea6b6ede03); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.ExecerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_defc836eab973764559fef7a6eacab29); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Pinger
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_f21346c5fd58f22896f02b5b5ed3475e); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_1bf71fb4be83d1fefe01cdd4d5a7a547); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_3ca02c0eda0e299af94cd79f7897d2b4); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_26864a2d0c8c5a605e497ec245ca9aa1); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_c4cbf403989171fe81ee42f1529c60fe); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_a1e981f752f5cdd201647cde399eb7ea); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Execer
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_70a606d71e13e67187acc5d5056da2e1); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_dea76879684fe95f8ca4bec7f7cec2d7); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_9743a8d2379220bf3bef64c2a633bd0d); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_40db605fa7c67cf6d53f67950cfe2a5a); ok {
		return struct {
			driver.Conn
			driver.Queryer
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_7848e2d0c0ff74bcf6995d5ad7cbbd81); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_b1c510487a94a905aff86496f31078e6); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_35700b3f7c770fcd0bf83836eee99fec); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_1c65fb5735328219b65d13d4e717b79f); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_45c298b388e6621ec2eb7a1ce46122bc); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_739d949041a357f14dd543ae1c82d4b4); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_7eb98af2c595d337fef40bcf00e367f5); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
			driver.Validator
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_734bbe2891841d705d0385cb4817474b); ok {
		return struct {
			driver.Conn
			driver.QueryerContext
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_a31656062b158c63465098f4075c1894); ok {
		return struct {
			driver.Conn
			driver.Execer
			driver.ExecerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_f2e0c829a13f0fc5e34c7487f6a5743b); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.ExecerContext
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_039f074dd8b7cd2fd02bcb0f3a89e29a); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
			driver.Pinger
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0002_004b3f5abaa12288a30861da40988df6); ok {
		return struct {
			driver.Conn
			driver.Pinger
			driver.SessionResetter
		}{c, c, c}
	}

	if _, ok := dc.(wrapConn0001_c3128ba459e963c7f07c0de8d58b775f); ok {
		return struct {
			driver.Conn
			driver.ConnPrepareContext
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_7a610c95e7c563485e18c228cc42e7fb); ok {
		return struct {
			driver.Conn
			driver.QueryerContext
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_f95e29e11298b900f8a05d72dcfe041b); ok {
		return struct {
			driver.Conn
			driver.Pinger
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_70145ca1951f1697e67aa095f55b8305); ok {
		return struct {
			driver.Conn
			driver.SessionResetter
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_dd0fe84666fad76cbf0dee440e953d06); ok {
		return struct {
			driver.Conn
			driver.NamedValueChecker
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_abc062388a1f7a907b94470050b8bc86); ok {
		return struct {
			driver.Conn
			driver.Validator
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_d815927e8e5ae23caed6edbf82ddb9a8); ok {
		return struct {
			driver.Conn
			driver.ConnBeginTx
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_f98304698b6793d84b0e7be539d135ce); ok {
		return struct {
			driver.Conn
			driver.ExecerContext
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_3380a887776f8c4760ee0d7b1fcdcde7); ok {
		return struct {
			driver.Conn
			driver.Execer
		}{c, c}
	}

	if _, ok := dc.(wrapConn0001_5f0836c749bf264b1429d45373cd34bc); ok {
		return struct {
			driver.Conn
			driver.Queryer
		}{c, c}
	}

	return c
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0010_77c7c6cd21c875211bd8f7fd0d4f7cac interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_c38f9867bd2650446eff3934abfe08ea interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_03ea25d7f45e22ec940d5495f1a84793 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0009_79f6c233903bd90e363a740e21dedd28 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_396e4ad5df4f8a9883b7db53e3150de5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0009_58e530768e77481673ca16ec9a1c4151 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_de85b18ab3bb83467b88d133c472e731 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_cf4b2a19600ddd6a199103b184a7b119 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_cc0840d8293f7ce49712448b9ab56686 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0009_3dd8e3401bce36381241be68e71a4596 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0009_d275e10eaed0aa246adacf06fbc7cb4d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0008_eb0574608d6a2bdeeda60227474f2696 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_b8a1712d42513721bb27fbf514e4dc08 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_7f39eb6d8011a9d621a372f56c9dbbb9 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_4efd319fba1b2f67427b7665301f5965 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_9741ae76b65cf39203ad2ebb1b042b26 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_58b49dcfceedf1a9157418d32ddc3566 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_5cf158fc22b2fa1ebca261e0519f686c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_5392d6d522ce153424f2d6a56e0b9412 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_a6448007f6fd3c5e82ac727e32dd9421 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_8a4cbbf06de72f30a062ff7d0d17a38f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_b62b5d9a91718c412fc82374a1ecc50d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_0ad7ad486ff238407782ecae921802ff interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_7030500b1e198b3921908510db90b60c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_fcea5f14de6e52ea3dcf885d9bbe6c6c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0008_77ff06f9507fa8e8907ea056fd8b3ff1 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0008_0af8ed4c453b8ed19240067f1ffc2691 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_92d776275e9ff04873b00304c46181c4 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0008_338d30b06f6f46413259d40fd65e6431 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_9424622cd873b695ec4a68179ae1c0d5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_f98a60ed56dc805bbb076e48569dfc95 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_638eb5dd5a685d42b2b21234d0f60533 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_3cca6f68d8b372d89bd75ae899a89fbb interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_93d6b726be36618017eb3d898d1485a9 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_79db6e70a15663317e248e989588499f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_83f7a44dbd5760c38c650a39e055e040 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_d5c1b115a4499cda31e8bc5417a3b20c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_d051dedc6dc0f7f58c5f991183648c05 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0008_519d23a1000bcef9e0c7ba7f65e2d8f4 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_a10c9902556dbe03fe0cb57fa734bae3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_d7069a0ed4e229445f0a98eb0a10e6d7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_31cab3e69b231de68861be7149ba4e29 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_a6e602e3b87f67896f9721167c748010 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_f10e8449922c36b75361a92bcd10a530 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_c5919040d602cb66dcdd424433bb7a3a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0008_3b2f99cd2da334bd09d06e4def61c4ca interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_5d41cafd5772738ac2731fcc2a5feac6 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_656773dbc70d65eacff8fb4b5f9e0aa0 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_4793858d8022e719f6eaf66c50c2d3ac interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_d45f7a4c9249a530fdbc798f8445be2b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_3d3aeb5bb32a489fc273e8ca839343c2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_fd07920044cbf1b6fcf3603ebb166c6c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_c35873a60d1144c83ef89ddfdd31f2cb interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0008_deb75cc63bfa66b8221455d806440ae2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0008_7dad59838d93adf969c56acaac791f07 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0008_615ccdcd04fa4562fdc78efaf9a3e52c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_448690263d84211bde0c1cea11a48b11 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_17ec4a11af4369f9cd9a79ffd22faaf2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0007_05d79f89f42e692717a39a29997da2de interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0007_9c83230a4aa511fa0866251a8cd1e0e5 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0007_52a0901ae15325b3f8f222a7da4e4491 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0007_7de64f0cac233a80cd163f2cea7cb779 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0007_267db5602cbea94318afb38924733c99 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0007_773881e5943867f63bfc186bfaf4c1c1 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0007_dcd5565feb96fad4810ddae84dbc06ac interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0007_9f05e5d4ed89fd9f8879183822aebe38 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_825b00edb1b43df6643231c638a67cfe interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_c600a1705c8565072162e90055bd7adb interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0007_854d6baadd0cebd58435d5cc29185337 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_13fbd5c85aa99a91119179870858fc36 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_d69fa1ea920d7f3f8644d49fc0b6aded interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0007_d9a28dc8353bb61b24670fc251b55173 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_41d773bbe4db96713017f5c8ef306195 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_d444a162b54c96890064f3460ddfa9e9 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_196ddefd25f93c78224271197dc93b05 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0007_f320700265733d40b644400d2731e9df interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_8a643673abce113f9fea5ce1702250a6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_bf737e0ca2c3b91f381fbcf33b58c0dc interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0007_c79fb6c89b83a91235030370f6e24c8b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0007_9274e2cae2bd02af387c9ce4dd6dfc10 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_427e4e33f79cba44d05a4674943c0473 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_01f5f56161c2ba3dc565775f6be2f8d4 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_fea3910637ab9db123d8a25f2431f80b interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0007_a7936c92dd92ff855159251fa022a982 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_5b4064ed641abac0e70005efb3a4ced2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0007_03a7f08f8d44561401428e04248ee99e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0007_265661e823b3aed3ebfb5ca65c6bba29 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0007_f321bc679577aebbd61ad4f183fdba3c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_b30963c7d0cef922992c78d9a6f4c315 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0007_68c6b44c1a5e01b29f32092c5f1d6cde interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_c74888bd57b65af1a21ac4d0e5030074 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0007_cf6f9229ac6b3ddb490a5b585f25091a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0007_41eaf4b14e8f3982b4305f1f8f903133 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0007_ed1ec44050b5b1aa1aff9aec0885cac5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0007_b8b951709b027931b2ae953b34b6805a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0007_1b0976ec19c297da5164b0249b0e47bc interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0007_ff474c2d7415a5c1c988cb0a1b03cd93 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_6d67a7a3477ea85d6035b87a3ae74f94 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_4b8fe7a4bef4b14dbe843b5a6815c43b interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_03145e837756150af12207ae91cc5699 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_dd667cce68b8575b61aaaa28533db684 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_68e4a2e940a99969f86eff9bf0200ae5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_e0d3b9c4b9a84b90126ca4e10d086879 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_08a849ae5b0aed116fc89ad126ee2ba1 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_2ced1cf7ed9e983a5dbdfcbd821607dc interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_b6cd19cf9a148786583dbf6fddc86891 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_deff6995f2d3779f620beab780c025cd interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_7e1a1c65bd708062952b2fa850109a58 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_9d60f26458c49e00bc5fd43d2603e619 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_3e05bf0f8875f261f4b42b7bf0b57666 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_051f43d212f7dad36785e63351b2e15c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_f52251e07be183ace8bcae49b3eb1f60 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_a0e97730c95bb2cd3448528f337ead58 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0007_e42bce9777e980fe284e9532ee788712 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_ba4fbb7dafdb5e52846909c46d3978db interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0007_f67e1bc8ad49300f830c4db205468d8a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_dd4cdd7fc99cf797611e88b4136ac1e2 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_05e0f0214217bf6dd30984948b2b6388 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_aba3918be493e52b0f4c7a39f4dca5a3 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_504dbaf2db1faa95f1f4b20785cd3e3b interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_71e6bf868d5a5f2785addbade2e3168d interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_40e6f60f31bc61ae5b41d15d9bedeb16 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_b46cb5ccfd4cf5e30bef58f218e7a7ab interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_9cd76498f7361882486b6d26fbbce31c interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_a9005e9e8fd249ed58a31bd4f3cefaa2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_01c3fe904d27d62954f4cc41a688e7e3 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_bff64378e85795f8fabccd7129574548 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0007_8bc59a3820b63b160785ea321b461a7f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_0621ad4aa976a74a406bd35b97479bba interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_6fd8c7ff8058ab175550eca5ed66d46c interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_0d92a7061fc421b6ab64b7ac074001ce interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_a715df8a59c6554c81c3f0a12913b128 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_e2480304f6e52f66531a98cd923184a7 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_c296c10f37d01b468fa5a6676f0ae7cf interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_ea59d342c0977cc9cff54e34cecd80e9 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_d672b036c3edb3bc6f356de7eb6fac7f interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_124ea9d194f98ad2d8ed38017c54db32 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_bad231076924f455873af75eb0b600cc interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_50aa24e8da7daf6237997b1479013bf4 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_6551ce77766b695dfca467061ed69cd3 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_42d966ceda8515bd592f89ae3ec5bf9a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_0991fbe5bdd96ef64ba2d9e4e193cc98 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_6f1ce97eab2efc4177f7a2e83aa9b03e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0007_6cc18f335b4f16ab71ea937da7bbeefe interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_f43f8c9b658de04e9329743f0314759a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_606569d605eb150571b585d33a7ac035 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_a2fb1ec32da13f4447d41d93a1c8f712 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_fc43713f9f34df269e4de565e839a72a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_70d68afe4880001b7820bce2c6562ec3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_a4184b99a265c878f5a78b178f8265e4 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_833f64043d5880be5449313155b52ad3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_0903cdaf440c013c89e41930fbe6e8e3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_e0cb9b3511ffadb5c06832aab335700d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0007_fc0e1e5d9ddbdc43483ef6559c753897 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_c44f50a0aa4c4ea4213e607c52ae032d interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0007_c55ff203396d2245b2978e1088c80597 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_eb39a4e3548820b6901cebed9b2cce6b interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_dff11d8a15f2ab6588886efe402c8501 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_8b66ea616b172e9ee49fd7acdf30b938 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_c9eac81b72ccb28661bbabb3f035f727 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0007_1210db1324cadd2baeae2e3872dae2c5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_c1afec2afbd894e5666ed8ff5cccddd2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_bdbf4466875acdb7dd6e329f0d22369a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0007_86034fd5bf9a19faed074c390256416d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_5a727ba6b819cdb39b9b7c0d12a10805 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0007_27c3ebf55b07bb44520d8c49cfc58f1c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0007_41d5d443d3bfb31a75284faff5774eb3 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_cba474b5c075f5f15ff51ac5f566807b interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_ba8eb78fcc8e48cdd4e64d9b9eeaa0ac interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_0e15d4fedf353ca2f5654696254bc8f6 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_9036e8ca98f6539a51570632bde39066 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_5ef510e1372bfdb56fe02f981a8624d8 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0007_ce0833b27dce8a567edf9a552aa55495 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0007_a9c99c789ed3a7ed62eb57e83d7ace3b interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0007_3642ed058ef985fd327d3616060baa99 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0007_6754fca0c4ee8359d1c34dd528128924 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_2787399cc3bf9f516640dc94361e2e52 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_4dd1be3407f9bbdb73f25370671091c6 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_5bbeeac10db407d89fbbbfe688174dab interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_7adcff1f6fdce3824e8a593ed868f703 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_36f3872f21fced8308f4e8b4173476ba interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_20a62cde5f5ee06f8316e1918b784015 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_6702d13c01fde7514f9cbb8fbf81348c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_cdf3b1387e4f674ad293aaa4273a23f1 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_e475241af1c8393981908620758d5e7c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_3539a7f57152c3ddf090b34d16228572 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_0e4f10ef28366c59fca887d987e448af interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_fd35073521bd29b9032b8e6343dc1582 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_0b9f1b7a8cddbcc442bda315e1f5aeac interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0006_08a845fcbe1fb6d2263f2a31c7fdef37 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_ab56b9fb37ba66c11b502a6aeb250a57 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_75a5bc2f988531d233f1385129322def interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_b82bb5e09eb92611efc96cd005703cd8 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_3dd755892540f437fc813bd4f9c6e76c interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_91a136fd0b4f47d3a4d4013f8095662c interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_1b8c4a926a314136f4f5b3211228eac9 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_8613dbc51a7af5956bbd224922fc8c97 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_8f10b08e9d138b1d7c3584d1ecc1451e interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_6f5a532ec7684dbd100cad131c0e6938 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_914f284da40915742abe2648432506ea interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_f1fa16fd964817de8b5d3eb8180bd696 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_d2ca153f50787fe8ff5d4d9d15df7387 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_600da2dbe24b60de75ed76e0e20f62d9 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_dbb137aa8bd155602797b2cf69146c34 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_4139e5345038cc13f9683ec9a592592c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_75faa058b6b2a57ebdaa36b174a9d821 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_8dc061cdbbf61be8784e703979185c75 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_8c2573bad4170e80835ea383cc59e9f7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_34d5ca3511ee070149519075ad8a1d35 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_2470eb657779092305966fce1e8bfcf2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_b1e6abfffd742d207c359fb991a8ff2e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_8e6b467944e3ea6eb4f2d78bfd8f4f92 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_d079c9555b86aa5a077b492b67851ebc interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_72a941a4cf0a009c149b56590a35c24d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_c7bf94d6c8d58e7532743914c4970352 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_ba4f7230813069263da91ea974ddb921 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0006_5b01e61de74885866a982713198a9e95 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_179505ccc3d31607ea1556bc528493f7 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_41fdbbe7cd36a1200433f6da51add330 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0006_6e8d1c43f28b1bc41b6db4315b0c53a8 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_9028f55e39a92542045cff44d23e5eb4 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_af905a4ef6af5f1b78600b6dc5f9e989 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_05e5e3b8142e65d950f76c265dff70a7 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_18c3ef0847a75edfaec330b73104621a interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_82c86c8ca99b8c772822016f8bcc6994 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0006_5d8d30cda81748d699df4c733cfeed7d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0006_9d9ba41da9501587188327a33b681b5c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_41877d4a96fe78c938c8dffea6a8a8d8 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0006_21c6aaf05ebd245d2a81458a7c024555 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_b715d13c93b8e3625ce37a07ca6653d0 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_f7af9ce566f7bb1681743a44b577aca3 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0006_0aaea350a70215c2fcedb8bc52c75009 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0006_92655dea42245236225225d1f8fb9ad6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_de66839b38d0d1407fc536262731eeed interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_63623688b7200b99d53ac1d8d149f446 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_57986a7bada864dc9bd27a63f6df85fa interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_ecac1f5922c3c6d122112bbcd2ea5417 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0006_daaf5cac24f21272a745e2120b9fd4f5 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0006_439a747f420d1b7d8502ab6eefe9a4e7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_792735c1a1541ec6975f0bc91c71e00d interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_47f531e2e2de8e5f660003b96048da4f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_606976a799355cbc2ef135455ae03a2e interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_b230083f64b1ffcc0c6bd334f364f4d2 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0006_0050743478792316eb42a9c00110d40f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_f18377cbcd9fdc9c278c7804f50192b4 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0006_9dfccfcd542b2f94cdbcd29ab48c84c2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_039f968629c44dd3bc1cb7ddff0d70f2 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0006_bfa31908831df498a7f6b53f63ad2cea interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_976d373277dceea1dca120722cddeb69 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_b546234b9cfccabcc950219c64db62a9 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_60cea439eb25c3dd4c5fc08377e33e86 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_4e47c679bef09f514d12456151e38722 interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0006_631a09a8b027fb27be1314d88c6c90d3 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0006_b1b8b2043027f1c9d648d17f54039e63 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0006_920646f9f9cd7de1ffc31a5841b4fbfd interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0006_3f35b9af9c8ffd4768e33f0c2f1a955d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_44acde8afd404ebe019eacfd337dcec2 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_db09fdc7c55f1435e80abe7fc49f0851 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0006_9aca8c5b7def495d806d3365865d62f5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_e6577aa3127a1ce75e3f57ec87f110e5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_960aac12246746bee40e9f7def10d535 interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_df7dc788652286f950681595791161bd interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_17f465977dd7e9292e49f3d3750b78be interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_4c7f8dd66b8f56bdc4c7d957a9350361 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_9e56d53a758766c7ef8bf18ffaf30cab interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_98d96fb40bc9afcd0640826a472020ba interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_aa58b7d423443aa26503c4665dc3743d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_4a0dfa7eb9abbb979297292dd3338e12 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_ef9a5be1db1823045e742d093b8e1679 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_f0703dc6977f55ca18452372104a81ed interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_c3214ef9cc08c38f4c20662dcae72311 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_f5575c07a41874f22f8f0b5ab4686728 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_c69289d8eb7ab88d6783ed0bae845fd2 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_069dfdaed1b51f9cdbe93d4f38bef2bc interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_bbd5fd8057adff78016752e1be0a9d10 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_55b10ecfd0dde9d7824b98feb70a4ae0 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_6ebbdc3a0d909a8ad477a64e2c788e4f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_ca53d6cee021350c5a67cf7963d3461e interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_33528021bc747485d09ebf30531c578a interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_1a5ff64aa5b26709b02c47092ec58a5d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_16ee6fccc64780a2fbe5db4b03c6d491 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_e25c3a362b3347dd39e21d3465ef0049 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_8a011e4074c65b06c74e7153bbde2b24 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0006_a41fec55a003e4dc7e9630066560fbd9 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_d6dd182e0deb72b02d08aca8d6474e39 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0006_f476294ed42363d7c2356114aab61ba3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_b9dfb6762a491c67ff005532fe4150a4 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0006_f2cd072da115c52699051ea5b4b4185d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_d50ee7efdb57dfc2c826d4e82c8b2f6b interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_45746a7b059c5779b0d5f2f07cc17264 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_82de9774803a1efd8ce92001a7d11f55 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_677628441b0b6583431e141a81e425a8 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_f412a098fffad4615edc59c07fb7df06 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_3eb7339761e80ca76a5e886ea74b5c62 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_e18b3805f6755a4369577290185cb406 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0006_bbc085109b2f0c1e0d02cf7fc8e946b3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_a0bfa772ad0992f8541f658be51b7492 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_560510e89015e29331831d9b1ceee341 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_f9ba701653a8c6eb07cf8297d1187aa0 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0006_ec4e435a7e759fb5447e7b6bad64a8bd interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_467ab682e7824a422bfeb422a169f5a3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0006_f3e8ae2b71ca5abd23475495944b773d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_08c836bfd3ac7defc7bfaa87a18a59ce interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0006_6aae8e835f06e746f9478a6ec320ac1a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_07b442328d533295abda3d01834c7881 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_1d1cc8f0334f9bf5c55f5fea90dbe9a6 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_0d1a47c831930910f2fdac0d91734bce interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_8ca52485803d8efbdfc5decc0c408e15 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_8f3fbdc92798a9fea1f7a0449454dd1a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_00546873757c60f5393ec281b0e23e2b interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_23e6077552c6edf7b4bf59bddcf4ed1a interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_7a5fa0f486ff80bb9c7fb0e446b45447 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_2c3fbfeef83ea24c40a3fb6dfa7d28cc interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0006_817a12cec18128ef9c600859b2480cea interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_2f219fe5f82e2f01a528964b1c0393e3 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0006_8153005fb418b0155b6e8959aba96479 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0006_f696e01b39fbfe1dd45edc9adeb4154a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0006_afb64ded6adc5fbc26c03214a30c979a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0006_6cad2465ee1d7bf5c6248515b4994362 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0006_2e099fd3ca48857fe0b9782e2c4be64d interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_7bc596b238c787350281f2424cc819a8 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_587fc8c53a8b9a50254f72cc34010eaf interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_480cee4bf05ec6c4b1cf7e4c4f5aa575 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_30c76ccd70969954c0ac7098fa6c4855 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_a78ba5b57311d51ec885c88fcf62a0d3 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_5802b2f9aac2139509d07fe2b7e4c96b interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_ec5f3f3ec9a278b63dda59c8c26ba82e interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_c9efb9ca182419988da01c441ce0020a interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_5046950504d8a3b58a4241a8d69a199f interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_70a752217140ab71d4273f5e892b9aec interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_0269e835f5f63b7464e26f2e5ced842e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_8830adf7455b83971bea3333d00d7582 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_fdb81facd138c4a0705bd634d6c5ec20 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_7ff2451bcca1dddd640a52a8a1d0b6e3 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_12a9fe7e3480d6e022d8342adce5e883 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0006_e6441cc632e3a36a18821ae55934537b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_8417a86915b2c41699e311e1f9d761ce interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_72f0714b315e98e5b6770a7be1a30036 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_06a0e4d619b983401d1cf6fd57b63b56 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_254aaf2cfda8272e56d4c15afba437a1 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_4186142f3f5124cf0aa4c9fcf856b2c2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_62f57ede9f4fdca4008bb8b102bff7ba interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0006_da7b76d23561eb1cb8672382675b31b7 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0006_71cea317b809e8f60b0684946b8caa5b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0006_61759473127617b1d4036288e4a8c0e6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_9e9e48a167e18445a9e409afe06df946 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_b023136fe1c19ee81c3975c762b0b8e8 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_6de7268e0b98d0ef7f3b4240d2068596 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0006_c5ce6bf14ca96cebd40fbe21bb29b7ae interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0006_8c4a4d51c5e2a5272c519b0c73f39980 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0006_da4c4979dc91a5536b453ae9de6ce031 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0006_e7b80b1f90f2a6e8ba69851be6bef3e4 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_346ac8866660b5adcd80229f1085a8e4 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0006_4ce176a6dbcbf31217f4f29b6694c68d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0006_99962fff9dfac7f527aec05bd41d626b interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0006_22ccdeeae81eec7d535dd685dbf5948f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0006_6bb429c97c0ecdbf2ad34c4ca1332c21 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0006_e771d868b1074b753fed0d5c9a3f9c23 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_59a3746398209af8d5abd1b9cc6e3ec8 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0006_728989255f27508682afd8f3bca8740e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0006_37e713c9732992f1fb1e77c0d97b5eaf interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0006_c36a0972036bee0218396f521005e737 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_3cd50bfb5d15730a5bde269d70e1bc1f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_904b3445479d8e81ff729ad84d948d4f interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_9a86c10ce63fdd3a02cb036abe9bf7f4 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_2c7c487b92de5c26457d4c323d171865 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_9bbef3126b2a214738ab649744f1c519 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_0de899839a124cfec12eb9923802d50d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0006_0ffce4a9e5d1475e6a58d770dc12865f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_aa8efeca905577f57d9c7b5b271a47c4 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0006_66fcfe34b767714399178ce8b69ff974 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_b1193c77c670778f320d60fdff28003c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_3c68d1110b4c053065a7d0b0d6bb71e5 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0006_9389c63162950678262bd1fce5c0b4cc interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_9c8b221b55cef0ba458722e6774a210d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0006_cbff1b1901fa4fcb5439f11ec0e5ab78 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_02ecffb5ba3ae606783c004a3b6d9580 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_8f4da85d1f4848b9967393bbbfcdbb5c interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0006_fc5da0946d81417cdde6ff2b10d5ec26 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_8f1d7cbe31c9ca537a8fb61a5f9f935f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0006_5882217b398b67e9b62b539def98ebc1 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0006_0021dd584f3a247e9c2e6708f1a62eac interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0006_47d7989aa77c482fa16a6787cc24626e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0006_ed3c586b0f2e0e36f908b9487bf6ec5d interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0006_7978a2c2f434e7dee7bf0f686b9a14da interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0006_beb0eb690c4037d495404a64cb1ff43b interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_2ebd13040db833ffd88ea5a8d2b377c7 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_cf0e2362c7b8c8b7b50e774f3e45a654 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_dab320a34a1bcc5b257599806c5751b1 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_5061a7994654594f8524aa2ca1c1d1cb interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.Validator
type wrapConn0005_eca72774318fc26016719a87b0035761 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_92148f5863b6f107ada35b49ff1c75be interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_8c8d9fd00cfea2ecf9be23c55594571d interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0005_4c04db0f86bdd7695eb22989598991ae interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0005_497c507b429fdc6223283c0a0a534ea5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_73f928359b60acd4f6a4d9ea2d69a0c4 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0005_777138835a2952632e61018f1c9e7854 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_3c4f23fa5c7397da192518024423d3ef interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_1df2cd699f63d177e5fb8d9548e2c0d6 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0005_bf70f73aa456d2cbf25751eb606ae32e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0005_a1fb2bb0c3e0b4e6e7a135ae626b0a33 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0005_b047b9570225719cb7f183bb00fa4daa interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0005_8cf3e7e6e07bad454a5b234c393a1d02 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_976d603d13330a49d799741fc8645f2c interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext
type wrapConn0005_6615ddf788dfa153e093b04b5c0304cf interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_5724728f6660a1e555fbcefe040a1726 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_b33d9f5f4ac39342663433fb14ff7314 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_65e1c77325da4d0c054e69bd4551c599 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_eb9908b466f20e92c27e6f5d7aa1bacb interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_951eb871b1cb468cf6c766c8c83e47e3 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_217bf44dbf7497d3f13ac203754a7e88 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_e3b29383c07170cb65cc52745a37ddc2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0005_c80cabe511a432e4d8e67766052875e6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_3da95f6ecd0aae49adefcccf5ee28c1c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_fe91be5d8911c57644753f925082c60b interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_513123ae151878de0171e3e4e1dc4c34 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_dab2c71afcbd6af7c7e4e507d436ebb3 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_0e4f2ce6ee8c6aaba65ea53cad0a7eda interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0005_d6b0f2f4ee5de04ac98eb3f61d53e1be interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_fec8999873ae214a7772c4e6cc77a241 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_d90d9f2acc455272bbb9d9676420675c interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_c05b0bf06ed751b21d3d661715644274 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_a1810061d6b6bd714973e83c4b126522 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_71ab8e9bc980bf0d31da9ff9c7c51ebf interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_8e89c1bd95e4e9731efc7ac1abba3639 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_7d6efd8ddf54ee184c44b622a1ee868a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_d6f7537a22515e5e9ca4f620125943e8 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_b09da8d8fd6084febcc163e9b40d8cdd interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0005_9126fe4aba9b3c897e19df3abbde4c2e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_48bc1b8fb59601a4f7543709bef7460b interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_3b53bb12eb2356991a95238fc2402972 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_95fa89461870d20a1aa400115550e236 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_370f70e1f7779d2b18c16b77796e6ad4 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0005_4a019bc6462dc411464060a64585b9f6 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0005_a5b461a6b03f2288d0023878bfc87868 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_61a248bb5bae07d449dffffcab033bd3 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0005_9e017ed93995c8a79a4cb81af48215bf interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_5e65f4368961d057065044067e0fdfa0 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_4d91aa8a5a995c9a419786a35a9f7da8 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_091322757a98b7a8c1de57f054e6eac4 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.Validator
type wrapConn0005_0b7241fbbde4c659f394f9d5233bfee9 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_f65dca199c8c7e749b3b1753cbfba37b interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_f8c02b123b7eb2a971596b935c942d7f interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_2e8be41ad679566886fbb5aa5471154b interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_650d3e7f112f289ca18b3e91964181c8 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_b1c239773538f5ab7b4ce96067f7a245 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0005_910e13cc5d980d54df51383c6deaaf0d interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_c4039cd70e55947c4038bbc29a3711d7 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_1b7dc36e8822afae1c8bb433c828d4e8 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_0e1cf0431fb2948288bf31c84c442195 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_db318a810481bb79ec5c2f495823168c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_2a7c43b8b90dba2d44286b98a9b73df6 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_e7dd37108246f62508c0e73239cf6e7f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_81570b25b723de4082fb023c40e4ccd9 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_436b3b6f6e649ed483a53237c3775c95 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_f2f5ce16d0204a3003a702dae993a40f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_d961dc8e20e0cf7d2f38aed674c88d64 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_633fbc849b41ff3a78e8bb75ce2db424 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_d37e6af15248fc04337ce571fb477662 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0005_bd2324d2dd9fc896512c6a5359795747 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0005_3c446fe7d7a33ac81540f3f0173a2617 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_8133b2ae04117f0b72b3d8ccbcc77744 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_7958686f0a2597a7cd7ec5775c25b60a interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_6a9a879c8cfd1e0a91ce0222094727c3 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0005_b442530703cde794bb889dcfd2337519 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0005_a47ece9818d8ba7d83c3dfee9cb5a5db interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0005_34bbe33cab14fe352dc7f15d5b8bd6c2 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_ff3f6f012645445f74f44e8d5a4f7688 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0005_8ea2641f961a5c26a90da3a80a453961 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.SessionResetter
type wrapConn0005_bb1dbd16911277023686c10027f764a6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0005_cded940add53ba9d752aac5b749677b5 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0005_afdced595d6aac26d33b2396172517d3 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Validator
type wrapConn0005_a7770d458fa6e73e3be636bc9da5a38e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0005_d2d66cb76facbf1d1b7a9c37bcf940c8 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0005_2b35a5af01307fa72cbc83a6a0a89f3c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_781918f6c2f2d7aae6eb6fb4b98b8ddd interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.SessionResetter
type wrapConn0005_1c6e0df7e9dde9612528f8d32ba1b892 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0005_26c5be906ea66f39030cc40db09bbbef interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0005_76190d9e70ad091b7248c8287842d7ac interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0005_63ba7a0c67d74cdff3b027a57f010878 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_39a8a13222fe1835ac5f0063826159aa interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0005_7081086f07db237298e2914bf8f98d08 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0005_55376846691a6cf458c674ca9055c9e8 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_f0c2880fecb021a3407210e94b5da7ce interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0005_dccec1fa0b1ac6076c36fab915923e7c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_20631a232115a4d613bbdd90bd786e91 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_db9df09b94b4eb586dda5bfa0be9970c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_76bf329e370e4f9a3dd99d83026a5478 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Validator
type wrapConn0005_1fa6f2b326e4409d7c306daf844e421e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0005_7a8098e8711012274339764352391092 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_95dabc404203cf349885a9d084b8cd6d interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_3dde8f11372ab3df00c06e9cc9f532a0 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_0c113831bfa3a7d8a21e204222a252c1 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0005_57e63123b079ce33af8bdadd71267625 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_43e95a1c134b51155bf9f0443e1d2244 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_e4c37a5e959f41b3c8f64b795b14239d interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Validator
type wrapConn0005_2b060163070164b37e556fa965a75b22 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_fbf26a0e2f1e9287b0f1b8a7673bb60b interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.SessionResetter
type wrapConn0005_2bfc2a716810b8b8e428184e9b039b54 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0005_7f384c69fc510e294062b2503fa418b8 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0005_15a9aa0a2fece20c5d53b7084ca5a21d interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0005_2d29a325eb959e4d792f3c431b67bd2a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_b23aaed3817fbc91a8f30fb2d067d608 interface {
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_c86504accefd86b5cd487e87c469c6c2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_f9d88387943c5f0ea9186ebb77de41cc interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_589a3cf2c8151915d84b023f02fc752d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext
type wrapConn0005_7536f936cf5efb906be9d5c9c4b9b797 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_d2d1a3dcf8158900936c7f1285b72256 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_43310fbccf884406bda0e2d6982345a5 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0005_9ad0458e4f0ef59ff7fa59f4c0837b04 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0005_95bd9777d6a5351f1331e1658e91c2c1 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0005_ca248c90d42e70e795469281c93508c6 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_399ec08db4e9d153838525614edae502 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0005_646f19e9526b1b663e651ef0e2bef0d1 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0005_7d9a114c2369646553f1a24e06faf975 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.SessionResetter|driver.Validator
type wrapConn0005_b486278e1b0983217ffb884258fdaa8a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext
type wrapConn0005_8b0d3838be5654907168a76bc7af037b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_be4447c4210552d0e2d709e04422e46c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_1c904363cb592579ae624ad523c9e4f5 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_c19f9e19c46b1102ea0299fc9a1b38a4 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_4fdf31a13c8ef55acf9bd6871ef8b7ea interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_0f04d6a673a007221d6a5b614f48b38d interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_e5655e77364e3f931a3ea6875c01a0e7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_24a1c25e2a27df99341328046a46e466 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_95c51c7c23affedd2c9b94e0b9376ad3 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_0ddcc5781c2afc1f94dc5b82a654a243 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_455bce9dcc99aa7c0fe96735f9c1c93c interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_76107da66f5dd1fc99f9f8706bdcbc77 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_13d3220ae21dcb3631a0e0d628adb44b interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0005_21c6c51021d4efd37fe596f5269fc7c7 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_f1dedecc1b8b0f34a2028626b0c5f9ed interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_2811e4468999fdb14616de0d05b9693d interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0005_61de1cea77f9a0a6f3a05cde41d82c0a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_bb0588c994e9ade1989ee58ccd0341e0 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_585bb225ed4b545c00cd962f5127ec53 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_50e95880bd580ebd349eb2d5d91ce2a1 interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_15b43c6aa30caf41244412fdba5c982c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0005_3221787875e443cec88a71604e26cbab interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_aa238e8cc38c595d62af9027950c4b6d interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0005_ac2671a78c195e70a534feb05d09ed8e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_2bbc56277fe4dd7194f1c426ceefcd73 interface {
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_036cfac2479cb17c204aeee52931f723 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_5b147a73931c3e3f8e956b0ee94f9d2b interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_a4f76c23639740e140aa2cdd5a506a3c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_fa8c6348b3aa58004e4cf20cfecc660f interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0005_cc1e37d25022b50f8aafb1efd4be0001 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_c670b45682d64d3a2f572ae08687b1c3 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_5698322bcb20d1da317fbe9951ea1ea8 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_be31557dfd123659214c16d5f27af1e3 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0005_c829a7b8ed999b05fe10801538e8fc57 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0005_5be0a90c30935b8c105dc54a4cd6009e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_b971b868aaf67fc950f128ca0831b10c interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer
type wrapConn0005_67ce16f2bfd9eae5a2cb635dbc252281 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_ebae4eedb7152c7969b661e0ef346b25 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0005_45e8e15a42cc8f475da57bb55786ce97 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_c159d592078beb8c2d756664e81eae2c interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_0a036aad328e938d09bf4e0bf0b1226a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_731476203a8004766e146f3ef7ded738 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_994f8417c053f4d92d2afd228765962c interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_11edbe1067aeaf4c860c11793036615a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_4309264d03baaeec83243ddba351287b interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0005_f86fa22f0a82b786389b11aa022dbc7a interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_ff3bb63d2da59b67dce62fc09a7068e4 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_271d69a36bda50544d2c8d10833b45a2 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_2a7c145d4ad450238d966cd1b81cda91 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_efd9a40fd7a5179c2d82d3a2361283ae interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_5f52d178514fd4128799774c187f794f interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_d10ae8f63596ae79b0e9b3e15c62b219 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_d8e90d85916412acfb7ae06a7cf3cd72 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_2a1026d851dd597c88f9d8d7bcd46f08 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_a613f56e071157bae4600ae04c11f63f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0005_726062f72cf8c22afcfef969d65edcd8 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_69e39a5cf47c21d10de2940a468498a4 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_c3e01796bb9f89d16ec7919b28547cff interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_d7c49a9c69464d5593c8d37de766beb0 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_2d5132626aee07585747967a992314ee interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_b9f4121827f5ff1843626bd5f67092d7 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_9445afc0bb6616b0ffb022c3955ca01a interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_ae364fa0447a19bedbfb479ca3ae419a interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer
type wrapConn0005_d18f395efe4408da08427fb8cf6b0f8e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_cd26b941286fbb8d42729f5e019f4e18 interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_a1bb57afc384646145f6f8acfff07ab6 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_34c6fd3f84dd7e5b474f79f57e7410d2 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_858ad5997c7a1ecd0c4e016da0b78cb3 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_451122d6894b9c1a4dc9870cc5c9f729 interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0005_5d88e7ee113b13cb24d3606d5b95737d interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer
type wrapConn0005_343af99016771b4bfeff7945efa2b0e0 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_3655878bb6cdfb4990bf98ac230ba6d6 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0005_9bbe32024047c4fedb1cec5108eaa882 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_612cbcd42b1ca988c32b19b0ee60fe8e interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_fd14a18f00f97dfb88beee48134b1ed7 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_f305594ed0d6b427fb2fe6e450424460 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_c391b2ea017883e989dd4c91421a2518 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0005_be38afe7bdce8819fe9f4816cfb05e2d interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.Execer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_fa236057bccd3e2547a2fecc91877c31 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_4b313c6d9a99a43b40d16ab65bd356a2 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_30e183bc3b73be79a53a3c5ccce7ceee interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_9820f7ab194cfc8681fa530ec810fbee interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_86594d133db62770dc866e4b90a41484 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0005_0ce7c7f80ec4cb8b2af6490440fb797b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_6f9178fa3ebe1aed2696357cd2afb3a9 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_1a8bfd6bc3b5d039e0d1ce8df766903c interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_5e64d450c4d7d27074726a3c99ca4467 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_a03e3132f8a1d5217b2ffa017ef0e471 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_a6894db2a6de9012f66db439c1729127 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_bf3dcedb268236b423e720a59c72d244 interface {
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_f4bb60d74c1105c14eef7ba20497175c interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_59952860d6c23bb78489fc55e73bc893 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger
type wrapConn0005_eb6a678f51255289dc34981b27eba51c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_063896be35045b544dfdc8bd6e379acb interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_522573917a875da9c8a8afd3f5bc8e45 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_86f6a80584134eec1e22e2bdf0b8b928 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_d1809e40a550459024766fd7df6b027e interface {
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_c19e3acd7931b2ef3bd181d17c473d6c interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_020e8fd0be2f1cfda622cd7fae226094 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_6134bd04f1e6222f7503813d44eae9a4 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_e1fb7970bf6eaaa98a1adf239f3ce9f0 interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_4f703912b1d28a906d4289a579b5e53a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_96e2182da02d44001714b6f444e74479 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_af6f1f46ce4181887d43e39ba0eac2a5 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_58c6e55019354e20f94dc97b5f01710e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger
type wrapConn0005_11ca78deb33bdf89e23765a32c4587c7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_20d8d229b2484368b368c970e6deaaec interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_172f2782b0ad97eecc7f56f17323847c interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_70c485a13a46b9c5b06e852ec8bb7c9a interface {
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_366a50e5d8552169108ca9bb15486b99 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_31d435075b7bc220179967c4a15393b4 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0005_f7594aa16f0a87f7e38788eb9c5f8e4f interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_2b9c34cec984ef83aa9d7eb139b39e1e interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_d98ff0a96153b4e2dd5d28a88b553717 interface {
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_1f8c595ee461d1900fce4cab06c9d1a3 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0005_d23443a383816cc688a8e378828508b9 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.SessionResetter
type wrapConn0005_6d4a973194a68c94feb8c4d04ea3b1e8 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_ba9fd7d66463d2daa1066b70ade39d1a interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_bdcb0b9943d4ec26d0fd0d86f47ba316 interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0005_af83d71bbe5f5d8c2f119446c9ee2404 interface {
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0005_c31c7d4b52033249a8a1d888dbd92671 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0005_91747f9e28ec013939bb19940a53742e interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_ec4d0075fd0ddc8b390599c00aa63ac1 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ExecerContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_a602fc8de20f9667dca43dab54118301 interface {
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_25d1043092b1d4fde69e802678864cbe interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_4abb5ab2b97be2656b1c86356a434af0 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.Execer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_c4837d1be7ab7d53ba06b807636c0d82 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_9e0ff90f7f6de1688dc2a4958037811a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_5497ec2107ee491c121c6b9957579ab8 interface {
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_7c2de3c593da1d0c0ca997bfb7e1d325 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_1b5b9005d979312b723d3a27c0c65db6 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_c130e5147dbf82d5c4a8d512eff1636c interface {
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_87e9112616a76f4e1f0250177acda175 interface {
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_09adbe5008d26d1ca9b6390889d3e56d interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_393cb9100ac7954d3e25afc8842b5dcb interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_bd27c65e18a44a6a66f0f10f280dd62a interface {
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_457719d02ae8d6b7caabc81e5ac31d57 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_6afb89b2559180fa651920ad74abbb9c interface {
	driver.ConnBeginTx
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_9efc0e3a9cf918d9d10523202afb880a interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0004_dcb06d9a6a9344136132f0325f437e0a interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_421d78dfe11851b6277447d6fc3d762b interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_62d2456de8e572dc055cc37196ab7b20 interface {
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_3b5220a9e11f03a530f5405a48a3583d interface {
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_b214fbc2704c5d1b1313515138517b76 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_e096a7508a321da18cfa19aba7f60b47 interface {
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_da3fc41871569e3343296bc91c49ae46 interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_62cc580731a221876d4d44f4e853a5e8 interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0004_54f6c56c9f04d106281a9ec1fa12677b interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_b67343921e6e00d89b71c1d31f11f704 interface {
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_2d60dbf9baa0aef32c08fa102e3e281f interface {
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_420b5a970fc5a50894e8f4bb8ee9ad8f interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0004_e3f29053589874b44fb4d5ed1443d8e6 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Pinger|driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0004_a5ce8739821cabdfe852f4ee9906463d interface {
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0004_6dcac6cce66b4cf60283b63f67f7ae42 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0004_c87c98d8350c1939120b60c35c42174f interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_7f966bf081d9facbed0cb718606514a8 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_42ee2395300cfbfe759565dd572359a5 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_7c9660a8ef7ac5680c6a0ec70aa74f7e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Queryer
type wrapConn0004_e3fef4b4dd1fbbec0e7621ad79e33759 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer
type wrapConn0004_41a10661ce28acc8f733fc3a3f82f02a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Queryer
type wrapConn0004_721d0a3d0378f788e179d749fd0700e5 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Queryer
type wrapConn0004_05be8063e1e61e6a640738667e26ff1e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.ExecerContext
type wrapConn0004_4a288634e9c34e5ea3325248b34c84d7 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0004_9f44333dbe686be726912364931d38f9 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_002ba8ae7ab8583e9fe9819c79f7f6e6 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_66435cfe8b3a54251d7169a90c4a43bd interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Validator
type wrapConn0004_575b7d6403f1859046727c7cc7fbb6ea interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0004_438ce4e1141c53dba704cd795cc91d37 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_ca46d67dab84b4a7356a56ad2734fb85 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0004_97fa4a641ec65c1c85b776a1ebeb05ea interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0004_3130001e0037ae9b22700b7a0c290575 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_c971b4273da1e06a47abc0e7f1d294db interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
}

// driver.Pinger|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_c9b76d5a245e1826763fdc264634a414 interface {
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_1b996e5d7c48c35626c06ddf2a0daee6 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Validator
type wrapConn0004_684b88bf9d678a38c20afbfe31f4589a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0004_cf59bff90cf956a08134977785d32526 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_bef873ac3553cf270cb99c91f8ad61c9 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0004_f9b0544f2b53b646873f2d45ad2a60e8 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0004_d72ac6dd8bc2e2fd59bef9bc1eeaf860 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_00169cf72ee7ba4bf707f0ee573f56b1 interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Validator
type wrapConn0004_31d4428785802ae57a327eebeda993d1 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.SessionResetter
type wrapConn0004_e0eb441d5b5198369a3d436dbdd96252 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.SessionResetter
type wrapConn0004_7b829e6cfa007bc3c4f5d8ee731ae57b interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0004_8497088c8479c1181a75be21c357c192 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.SessionResetter
type wrapConn0004_8e129f1c2290819ce637f1864c5db12f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.SessionResetter
}

// driver.ExecerContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_1492bd1459024954efa63e5096a5ee23 interface {
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_4d977f1346507e064ab5be6f55e3c261 interface {
	driver.Execer
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_0b8b629d5421fa2672c401e0c3a208da interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Validator
type wrapConn0004_86d48c27c6f553db31d7023a0618aeab interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Queryer
type wrapConn0004_1f7dbf76ba241759ea22fd9e67b0c97c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Queryer
type wrapConn0004_4dac992fd9f94a93dffaf95424246fa6 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Queryer
type wrapConn0004_6a4c903cebec66621da83b719ab797e3 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.Validator
type wrapConn0004_8c1f1ab2627a320478eb2bd9c49efd73 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0004_f679a21ffd1bb700d447688e1f22b988 interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.Validator
type wrapConn0004_3e0d6a60d0b94f9a38bd6278a9aca508 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0004_1d4652a872051fa3240ca534bef0a43e interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_994bda6d7619bde46104d2eab7ff768c interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0004_63ada1d75739c374aa86100675e5d5d1 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0004_052895537635d2f5a24651bcc48ec975 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0004_ecaf12077437db7ea718765cbeceaa78 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0004_8fb2e6de52e01bddf84ba9aa01f20489 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0004_c6e51aefd95b2286637499cab5834392 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0004_6e5b2720e60eff99a0177f056c13773b interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0004_e480d4975fde5e93e67c0901d131c42e interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0004_7bced5185e944b855170b3b92901575a interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker
type wrapConn0004_7e1db962b81a09f792909878ea32a5fd interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0004_7132619a8c83606d629fab216a12c395 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0004_4272cca8760969ed2322624711f94407 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0004_8098fe71a9a9219d58e590ca89192425 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0004_5d17fc1b7064520ae6022fd2f62c2bfb interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0004_6fcd36e5fdd1db57ef48d5bdd4beebd9 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0004_ff336118365ff6fcd2a2e761486717c7 interface {
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.Validator
type wrapConn0004_62deeb4f7a0c751fdee63cefb834691f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.SessionResetter
type wrapConn0004_855ad8d8a1e77072d4e79e775df962d8 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.SessionResetter
type wrapConn0004_59eff3dd52fe69e24ff9a0e5fbfff105 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.SessionResetter
}

// driver.Execer|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0004_ee03ba98204ba8f2d9fbb2e14d6a728c interface {
	driver.Execer
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.SessionResetter
type wrapConn0004_8fef826c40cffbdc6a69c70b5676c824 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0004_62083746252ebfd031b91ac334e101f1 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0004_2dfdae3f68bd27ada2cc3a92559282e7 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0004_c68aa300799a706219ed2f7ad15362c5 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0004_c948501ec7e870a27ec7845843aeaabb interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.Validator
type wrapConn0004_68b9285a81940136122cef1c146a64b8 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.Validator
type wrapConn0004_0b99abaa61e35521eb15190ed9f7df0e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_81ab8fa6399bc4cb3f8b4ceef191fbe1 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0004_c44f3469765232bf2f3748ac3491ee68 interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.SessionResetter|driver.Validator
type wrapConn0004_ab55d13e991c298837fd0697904e660e interface {
	driver.ConnBeginTx
	driver.Execer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0004_573d41654139e877dea050658666ffa7 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.SessionResetter
type wrapConn0004_b560d15b65f988c4960fc412191c25df interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0004_72837a807cd7df18bbd7e22ee697c816 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_abcb14740cbc0d1a201f9176fb24d825 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_9adf723bcac874299a491196558e7dd7 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0004_5f6080df6bedd5e86600c0719613d97a interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_4f64f93a4e925dbd2a7598b37bb8f332 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0004_60f92c1f81ec8480346add9207661f84 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0004_f33d7e93149b996e3b7d599a621eb41b interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0004_8480b7a2f4b03ada42be3b2e1e8ee03b interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0004_2a24acd339cf38eb1b877ff6e36477c2 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.SessionResetter
type wrapConn0004_b702ceb57e9a027bb97820e611754bef interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
}

// driver.Execer|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_be26b6411cff1ae6721379214d087f9e interface {
	driver.Execer
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_daf8acd200daae2412c1313c41d2df26 interface {
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.SessionResetter
type wrapConn0004_089675b7ca4c9fa833b6515226a0b60c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_b454cff4daf9ab30e3870da405d5fc28 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_86f21eab726445da5b1ba4da19a01e24 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0004_b0e5c7fb09f6acbdfcf6e0d4258d4a9b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.SessionResetter
type wrapConn0004_ad424cf88ec9bc306be6b23b4034be4b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.SessionResetter
}

// driver.Pinger|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_ff5d54172f489033d905c69de9288f35 interface {
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.Pinger
type wrapConn0004_86340b112896f6992ac55fb1bf8e74b0 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.Pinger
type wrapConn0004_951e06b10a39de1cad7dccca3392ea82 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.QueryerContext
type wrapConn0004_74d560cdbdf369e005ce1b8ce93a13a0 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0004_a6112ca861997f87f6d0e56ce5b1c43f interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0004_a90a33e85b58caec9ced942971c780f5 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.Queryer|driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0004_1bdf28e924b92fb4795814702aa75ef8 interface {
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0004_9050482dfdb75ed3bda0f45ad2cc37e1 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext
type wrapConn0004_65ab36fec2bb30a1596eb25d39a63d3c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0004_15e77a0e9aeb76ab2f04aa748ac4719e interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0004_c8f1b5b5a8743b962a62682d4a6e9262 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0004_8cc60054c0093e34e396607e0b128191 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0004_0fd04b4e30706876c57056f2a95c5b01 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext|driver.QueryerContext
type wrapConn0004_1915c1c223496f64670367e8aa3e3eed interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext|driver.QueryerContext
type wrapConn0004_554f29b91fc8ef3f3adbf7b1224f0a44 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger
type wrapConn0004_5423a435d10b320ffac1aac27fe43c6f interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.Validator
type wrapConn0004_ac5d94b9c5d9718d65ebc6e6644466aa interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_d3d3167450d7d79e6e587b43b5014c17 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.Validator
type wrapConn0004_4073315e2d95994c6f3228bd979f95d6 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.Validator
type wrapConn0004_270c1587b679b20fe7add8f451fca302 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_51b8c82cc36cae2aa27ea9f5828ed420 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0004_a02b8d7f6b52abd4bf35c4c2828fa170 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0004_2e864f41a5216ee217adff2fac1e6836 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0004_211f8a32f588ff83dc85e793d7ea741a interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0004_4adfd3cf6b9fb282abf8009ccdfc1dce interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_40b177ca84280f73c468ff30cdd2e97a interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0004_e03d707ef6c5cc14794c82dc511c6c91 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0004_1d2ef1d66fc107f4d6470ccf34b8b037 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0004_3fc6609110f1567ea6fb2ccc9c6941ca interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0004_a93de0616ed3079543f82d35645aaff5 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0004_aec80ba98ebb4fb6dfb693519c486fd6 interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.SessionResetter|driver.Validator
type wrapConn0004_47a3327409e70c795fa50aa402f552f8 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0004_020f136197f532507b2b739e893298bd interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0004_fbe9055480e6ebbe349358be151b7b87 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0004_2e6272eac15fc86b3a447774b973a7db interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.Execer|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0004_fed4efdfca6c7311372ae1d21591bc51 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0004_eeade78a4548ef600a551ccc3677a14a interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0004_8b02254c6cfcd13726ae202cdb121b1f interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_f018b29703e3acb35c78e692c0b4e037 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ExecerContext|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0004_537d8d4f3430c6871d023b75293bca4d interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0004_c6260f7de558b080743c3c6705aa04c6 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0004_1666fda37418e8810185898fd16573d2 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0004_a0bae6dc5e6fadde4228e11bf5e49dd1 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0004_bf56dd41a04eebaa9037e34d998af495 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.QueryerContext|driver.Validator
type wrapConn0004_4501f160909b47a29559b718409520ea interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.SessionResetter|driver.Validator
type wrapConn0004_4eafb813d5396d818b628d6737e80d42 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0004_4183310106ce1532189899565b015546 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.QueryerContext|driver.Validator
type wrapConn0004_e741bd5513728d7fdcc7ea4757864d07 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.QueryerContext|driver.Validator
type wrapConn0004_bf95b1c0c61b8c9f9cb99f3e15ce6042 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext
type wrapConn0004_28c597b75ba228e5d6902ddc68ad7811 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0004_5de483c883d2daf2738219ea08483e30 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0004_4ff0f4092dd79c548ab0e495e6cddc3b interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0004_c4754bae9bc3da992a48068120977d41 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0004_466daa400cd4119b35117e42b62f7a5d interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger|driver.QueryerContext
type wrapConn0004_819ead64aeb08e11d926e74832dc7856 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0004_4ef020c4936aea427d194c2b96378df5 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0004_12ebe2850e61f4a8e07fae374063118d interface {
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger|driver.QueryerContext
type wrapConn0004_30facce362b8d7c83804cc9435f8bf2d interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0004_a0f6a4ff15d8d9e4f751f67c2efcd8f6 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0004_c818cf1c81182b6fe0e9f65adf919fcc interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0004_94eed8c577d89cc2d2f123f92b3e0a4a interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0004_ac82cf73c860771ac75efc0b39b24537 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0004_d3598002e745fefcab218230adaf46d1 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0004_86042fa2173f60c1aaa263c2c4f6b7ee interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer|driver.QueryerContext
type wrapConn0004_28f0c23a935aa8023225aad3b1586102 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer|driver.QueryerContext
type wrapConn0004_418e5123313e62fb4031520b233885ec interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0004_47c38651973c418d8726100d2bdf9d3a interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext
type wrapConn0004_789e9da15950f1b6a8a6c4eb8902709b interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
}

// driver.Pinger|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_5f52c927d92988b4856ae8f434c23f6e interface {
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0004_c921c04aa70e11394e15484e40785309 interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer|driver.Pinger
type wrapConn0004_1df93e659bfdd7643b71682374b5ee3c interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
}

// driver.Execer|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0004_c7a900ae364616340bb424043839d4b2 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0004_645074a289d63811b9e9a8978d5ba260 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.NamedValueChecker|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_249f8105c1428584ff3103b40440952e interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0004_555cdb27b272a3f2c0eb9f7b9e0eb5ff interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0004_d1052a967ff8db9233909609d294e184 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0004_c5f13f7f8f3b7d0270ae28c5a1d96ba9 interface {
	driver.Execer
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0004_51825a78fc953cc98d20358550621612 interface {
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0004_f39e7f23bc4a7c1bd04d47bcdf9f2b75 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.NamedValueChecker|driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0004_95c3412bca03d21c7267d3279e2f3312 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ExecerContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_42ed5cc338556317302dcabae95aecba interface {
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_66f15c6b6bd4dc2a63272af2d58fc85d interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0004_9fc12f8e3f7d00cb13a0a0b563ab71f7 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0004_af5eb26d3936b8bf476e1ea60b3edaac interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0004_84343f2e4ec9d07e1a9dbe14088de505 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker|driver.Validator
type wrapConn0004_6d1b579154234b2d3345d0013c542993 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0004_06ee4205d01722f7549f02e0b7087b0b interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.QueryerContext
type wrapConn0003_55b0c9c15fb89c88e545e8b9afe973ee interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.QueryerContext
}

// driver.Queryer|driver.QueryerContext|driver.Validator
type wrapConn0003_df39e5dbce330518630020ccd9345df3 interface {
	driver.Queryer
	driver.QueryerContext
	driver.Validator
}

// driver.Pinger|driver.QueryerContext|driver.Validator
type wrapConn0003_ee7a3b8aa7f629ea0410d3cb5df23364 interface {
	driver.Pinger
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.Pinger
type wrapConn0003_17f4ce77ecc13085127e932ba6a8fb65 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Pinger
}

// driver.ConnPrepareContext|driver.Queryer|driver.QueryerContext
type wrapConn0003_f168d215225c7302da641eb4505c6093 interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.QueryerContext
}

// driver.Execer|driver.Queryer|driver.QueryerContext
type wrapConn0003_a5e837eca1906a3008beb989e61b545e interface {
	driver.Execer
	driver.Queryer
	driver.QueryerContext
}

// driver.ExecerContext|driver.Queryer|driver.QueryerContext
type wrapConn0003_3ba01671ca6644b664bff3737b516c82 interface {
	driver.ExecerContext
	driver.Queryer
	driver.QueryerContext
}

// driver.ExecerContext|driver.Pinger|driver.QueryerContext
type wrapConn0003_80b435a335887493ba9656f5a0e76a14 interface {
	driver.ExecerContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Pinger
type wrapConn0003_1bf343887515a244e3d82298438693bd interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Pinger
}

// driver.NamedValueChecker|driver.QueryerContext|driver.Validator
type wrapConn0003_0040353f748f8b98cb8da172c6b336c1 interface {
	driver.NamedValueChecker
	driver.QueryerContext
	driver.Validator
}

// driver.NamedValueChecker|driver.Queryer|driver.QueryerContext
type wrapConn0003_714551a2b9e17d67c6952821e0182bf6 interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.QueryerContext
}

// driver.Execer|driver.Pinger|driver.QueryerContext
type wrapConn0003_e00c84ffbe65ef6bc98ba49b16f282d5 interface {
	driver.Execer
	driver.Pinger
	driver.QueryerContext
}

// driver.ExecerContext|driver.QueryerContext|driver.Validator
type wrapConn0003_086984675e3daca69cb4d93a08b03c05 interface {
	driver.ExecerContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Pinger|driver.QueryerContext
type wrapConn0003_ecd209322f30837e1dde5cc3c5903232 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Pinger|driver.QueryerContext
type wrapConn0003_f589cedd98232e73dba33f2b1ee57f05 interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Pinger
type wrapConn0003_1c8897aa6698e41ba80853991cc24839 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.SessionResetter|driver.Validator
type wrapConn0003_7908fe1d05a5c3552e8ea5c1ee7d48e7 interface {
	driver.ConnBeginTx
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.SessionResetter|driver.Validator
type wrapConn0003_895867613f93da07976594c2ea883249 interface {
	driver.ConnPrepareContext
	driver.SessionResetter
	driver.Validator
}

// driver.Execer|driver.QueryerContext|driver.Validator
type wrapConn0003_739075d41fcec00bc42fc215117d32ed interface {
	driver.Execer
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.SessionResetter|driver.Validator
type wrapConn0003_e955cc61481c7908edd1588edb5952a7 interface {
	driver.Execer
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.QueryerContext|driver.Validator
type wrapConn0003_6168ba655398488893634accebfceca9 interface {
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.Validator
}

// driver.ConnBeginTx|driver.QueryerContext|driver.Validator
type wrapConn0003_975ad90f36f2a8d5d14cbc769a6ce607 interface {
	driver.ConnBeginTx
	driver.QueryerContext
	driver.Validator
}

// driver.Pinger|driver.Queryer|driver.QueryerContext
type wrapConn0003_46d10e7d36cfa34b2e97b70bd8f4e8d8 interface {
	driver.Pinger
	driver.Queryer
	driver.QueryerContext
}

// driver.ExecerContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0003_42faf5a4eb464290a69fa9800a551bd4 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.Pinger|driver.Queryer|driver.Validator
type wrapConn0003_ca97cadb92a7dfcaebf19bed8659a2d0 interface {
	driver.Pinger
	driver.Queryer
	driver.Validator
}

// driver.NamedValueChecker|driver.Queryer|driver.Validator
type wrapConn0003_fd6abc6efea6a6cda4203b6ca3a23574 interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.Validator
}

// driver.ExecerContext|driver.Queryer|driver.Validator
type wrapConn0003_b0d6367c8ae272dd20ac1a3fc8269d5e interface {
	driver.ExecerContext
	driver.Queryer
	driver.Validator
}

// driver.Execer|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0003_f71834a64db6e15385899741cd4f16f7 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.Execer|driver.Queryer|driver.Validator
type wrapConn0003_d2dda163c53609ef7d4070cf89140c06 interface {
	driver.Execer
	driver.Queryer
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0003_4a1ee3cec4754da65ab631f8d55709d8 interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.QueryerContext
type wrapConn0003_86f0c8258522256415a6791126d2786c interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.Queryer|driver.Validator
type wrapConn0003_89bc20456d7b564e816ea43252326002 interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.Validator
}

// driver.ExecerContext|driver.SessionResetter|driver.Validator
type wrapConn0003_065f30b751342fccddaa383211967ee5 interface {
	driver.ExecerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.Queryer|driver.Validator
type wrapConn0003_d9070b80f72659573a9a1e1ea8b7fe38 interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Pinger
type wrapConn0003_4afcf6b296cf46480b7f20345361eec9 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Pinger
}

// driver.Execer|driver.ExecerContext|driver.QueryerContext
type wrapConn0003_4a6b7427618f96e139a937955fa945a1 interface {
	driver.Execer
	driver.ExecerContext
	driver.QueryerContext
}

// driver.NamedValueChecker|driver.Pinger|driver.Validator
type wrapConn0003_4e478d169a0c925294ee737c49c1ac6a interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.QueryerContext
type wrapConn0003_58164b71668ef615a11f9d6dd3ff1e07 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.ExecerContext|driver.QueryerContext
type wrapConn0003_66ae780e5764601ed26339f0ef146ab9 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.Pinger
type wrapConn0003_3e97f3c285abb15a2fe7d2b9982ba7cf interface {
	driver.Execer
	driver.ExecerContext
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.SessionResetter
type wrapConn0003_c02e55bb6f1875e5ca0c1667579be31a interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.QueryerContext
type wrapConn0003_e1b111d2c8c09cfe4d5b9d367c734191 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.QueryerContext
type wrapConn0003_8ac86baaa7af7d41aebe48b4e0583658 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.QueryerContext
}

// driver.Execer|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0003_576714a6f8eb9ab037040fe7ef742b28 interface {
	driver.Execer
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.QueryerContext
type wrapConn0003_cdbed2933e5c0be5f6806d38207387b1 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer|driver.SessionResetter
type wrapConn0003_21f118405ed242fbb797757d53cf8761 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer|driver.SessionResetter
type wrapConn0003_4aab8206dcc15dbead696d4dd7e33cea interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Pinger
type wrapConn0003_5a07d33df77df1576ce8067b8997d413 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ExecerContext|driver.SessionResetter
type wrapConn0003_348becb0e4f2e9da375f40a1e59807d1 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.SessionResetter
type wrapConn0003_429af3d54eff18d5961d7d87b53413b7 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.SessionResetter
type wrapConn0003_9b9b5f504f5a2fd077ca61e40c8f53f1 interface {
	driver.Execer
	driver.ExecerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.SessionResetter|driver.Validator
type wrapConn0003_f3e446fd938f183e46a975476439f43d interface {
	driver.NamedValueChecker
	driver.SessionResetter
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0003_ff274753683e1d0adf6ecc55b9be8631 interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0003_67a6f30839ce48a2e9f3488bf5e2053d interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0003_4b5d9a137b792a784a0b494e06773edb interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0003_b493216de6d7d70b1fa112fe946c8bc0 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.SessionResetter
type wrapConn0003_2c310269b09748f7be39d17b8528f2df interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.Pinger|driver.SessionResetter|driver.Validator
type wrapConn0003_d033831362fd7cd8abd09076a328bf5f interface {
	driver.Pinger
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Validator
type wrapConn0003_57e946d049a38df9bcfe6ef5d8598cfc interface {
	driver.ExecerContext
	driver.Pinger
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.Validator
type wrapConn0003_1ce67851a70243fd021fc5ea522ae189 interface {
	driver.Execer
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext|driver.NamedValueChecker
type wrapConn0003_ddd509d4dc401e51bac497647969675e interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.Pinger|driver.SessionResetter
type wrapConn0003_ce39b6405b9206d7ba9fe332b19cf89f interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Pinger|driver.SessionResetter
type wrapConn0003_66e6b00152e5329959d8d66776407655 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.Pinger|driver.Queryer
type wrapConn0003_95d38f7324a5d842e83132f8f4476f2c interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.Queryer
}

// driver.Execer|driver.Pinger|driver.SessionResetter
type wrapConn0003_ea65c7fb86fc48081e1d867a9cfc1870 interface {
	driver.Execer
	driver.Pinger
	driver.SessionResetter
}

// driver.ExecerContext|driver.Pinger|driver.SessionResetter
type wrapConn0003_3371fc70ba5d416a901c5118687e778e interface {
	driver.ExecerContext
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Pinger|driver.Validator
type wrapConn0003_aebab2fa61e210e3c96cb924b4356379 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Validator
}

// driver.ConnBeginTx|driver.Pinger|driver.Validator
type wrapConn0003_b2d88477280fe09a57ad676180ae7fca interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Validator
}

// driver.Queryer|driver.SessionResetter|driver.Validator
type wrapConn0003_e1f4671a79f9b4bc9f39c5b15d6aeef8 interface {
	driver.Queryer
	driver.SessionResetter
	driver.Validator
}

// driver.NamedValueChecker|driver.Pinger|driver.SessionResetter
type wrapConn0003_75415c9d129d1dd60031d39ea83f5316 interface {
	driver.NamedValueChecker
	driver.Pinger
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Validator
type wrapConn0003_344b3b12f6cdf795bbc480ca3c22d0eb interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ExecerContext|driver.Pinger|driver.Queryer
type wrapConn0003_907ef3b14632933f40c5231177e7df3a interface {
	driver.ExecerContext
	driver.Pinger
	driver.Queryer
}

// driver.Execer|driver.NamedValueChecker|driver.Validator
type wrapConn0003_2295005ac105ed3d2e96ca764c9284e8 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Validator
type wrapConn0003_55276ca2338c0cd130c39fb53ac4caaa interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Validator
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Validator
type wrapConn0003_ddca4fe2167f6deb65f9ee0d4d46b23b interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Validator
}

// driver.Execer|driver.Pinger|driver.Queryer
type wrapConn0003_24b8a6effd27e68d07d33d0e6d4e7124 interface {
	driver.Execer
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Execer|driver.NamedValueChecker
type wrapConn0003_01d353c46df9c41bbdd5f105418d9b12 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.NamedValueChecker
}

// driver.ConnPrepareContext|driver.Pinger|driver.Queryer
type wrapConn0003_26a94233bac30452998cfcff3d4ab4e6 interface {
	driver.ConnPrepareContext
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.Pinger|driver.Queryer
type wrapConn0003_5701107029804f6e326b6191c81f660d interface {
	driver.ConnBeginTx
	driver.Pinger
	driver.Queryer
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0003_e8d66097d62d92ca5e8a01d0627ef47b interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.Execer|driver.NamedValueChecker
type wrapConn0003_d850e679e03edc84e19f7ec96957d9bc interface {
	driver.ConnBeginTx
	driver.Execer
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.Queryer|driver.SessionResetter
type wrapConn0003_cd50db51f3efd4e06b314497df4e8bd1 interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Queryer|driver.SessionResetter
type wrapConn0003_14a98e8d04e2b4b3493f4673295b2dcb interface {
	driver.ConnPrepareContext
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.Queryer|driver.SessionResetter
type wrapConn0003_c09291a9dc155f8d0c0f1d1961c478c7 interface {
	driver.Execer
	driver.Queryer
	driver.SessionResetter
}

// driver.ExecerContext|driver.Queryer|driver.SessionResetter
type wrapConn0003_87c58cadd923ee575179af1ec5f175f5 interface {
	driver.ExecerContext
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext|driver.Validator
type wrapConn0003_251125afd2a9bfc7c7bcd717ad73c7e4 interface {
	driver.Execer
	driver.ExecerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Validator
type wrapConn0003_3046c3f6840014e9d6851479421aa2db interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Validator
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0003_0a04ef57a9627c79f86b26469450314a interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Validator
type wrapConn0003_214d339ea349168fbe3e0955424c5eef interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Validator
}

// driver.NamedValueChecker|driver.Queryer|driver.SessionResetter
type wrapConn0003_84c6caa443f9cbb04c70721344581649 interface {
	driver.NamedValueChecker
	driver.Queryer
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.NamedValueChecker
type wrapConn0003_1a34bf796737439096bfa6f10d5696c2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.NamedValueChecker
}

// driver.Execer|driver.NamedValueChecker|driver.Queryer
type wrapConn0003_35e2a252407d971a4249ddd0ab337396 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Execer|driver.Validator
type wrapConn0003_47814f9637c1fd14af8c27076a6fc4be interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker|driver.Queryer
type wrapConn0003_a35f893b3d8e8eb80b18a381c5ce721a interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.NamedValueChecker|driver.Queryer
type wrapConn0003_485292d166cf8930ff7f1c6d74199750 interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.Validator
type wrapConn0003_e575d519c87b84d55ffacfb52b72183c interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Validator
}

// driver.QueryerContext|driver.SessionResetter|driver.Validator
type wrapConn0003_613592023fc524d8f895420661ca88e0 interface {
	driver.QueryerContext
	driver.SessionResetter
	driver.Validator
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Validator
type wrapConn0003_ecd9017910d8eaafad4c84916460cde2 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.Execer|driver.ExecerContext
type wrapConn0003_4b7e2498575c5cb9db1b7ab70c55ed6e interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.ExecerContext
}

// driver.Execer|driver.ExecerContext|driver.Queryer
type wrapConn0003_dc6b33036e44bc5045ed542d781b64c6 interface {
	driver.Execer
	driver.ExecerContext
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.ExecerContext
type wrapConn0003_4b20cfcfb6db2db4449ec9444be99683 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.ExecerContext
}

// driver.ConnPrepareContext|driver.ExecerContext|driver.Queryer
type wrapConn0003_763309d9ccc50f11156cd9497d38589a interface {
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.Queryer
}

// driver.ConnBeginTx|driver.ExecerContext|driver.Queryer
type wrapConn0003_054b90a03b02ee83d5b32a2618de1414 interface {
	driver.ConnBeginTx
	driver.ExecerContext
	driver.Queryer
}

// driver.Execer|driver.NamedValueChecker|driver.Pinger
type wrapConn0003_b1e4b1e4747a4180955083debd843744 interface {
	driver.Execer
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.ExecerContext
type wrapConn0003_0980f76b2d444d923da26fc4a1b82818 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
}

// driver.ConnPrepareContext|driver.Execer|driver.Queryer
type wrapConn0003_a3a706415d197325957442aa15129e61 interface {
	driver.ConnPrepareContext
	driver.Execer
	driver.Queryer
}

// driver.ConnBeginTx|driver.Execer|driver.Queryer
type wrapConn0003_1e0e8df2036eb4f1f761196bb7b1da89 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Execer
type wrapConn0003_b9c9b9171bfa8c16742e5f91689b5d88 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Execer
}

// driver.ConnBeginTx|driver.ConnPrepareContext|driver.Queryer
type wrapConn0003_1b7e307d8021743b2d0cec1048d55722 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.Queryer
}

// driver.Pinger|driver.Queryer|driver.SessionResetter
type wrapConn0003_b49194b335bfb08661e6e40529986bad interface {
	driver.Pinger
	driver.Queryer
	driver.SessionResetter
}

// driver.Queryer|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_35f2471903ecb7c1e4cde7fd5cfe5e1f interface {
	driver.Queryer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.NamedValueChecker|driver.Pinger
type wrapConn0003_68f463fe3d9a985f2c57f5f872c7c333 interface {
	driver.ExecerContext
	driver.NamedValueChecker
	driver.Pinger
}

// driver.Pinger|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_d40c325b5125e64636b9b8ea3f1b8975 interface {
	driver.Pinger
	driver.QueryerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_9d0152a9df319a9e7e4a0bdc7f7fa2d4 interface {
	driver.NamedValueChecker
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ExecerContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_c18fe73278d7622bc6a1b7ff02b795e7 interface {
	driver.ExecerContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_8654b308d9ff695cc51971a84d993c51 interface {
	driver.Execer
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_67e361827a469e5a7fd41b0778aa1958 interface {
	driver.ConnPrepareContext
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.QueryerContext|driver.SessionResetter
type wrapConn0003_23392d2c642abd59125f77ecca6f7b8f interface {
	driver.ConnBeginTx
	driver.QueryerContext
	driver.SessionResetter
}

// driver.ConnBeginTx|driver.Execer|driver.Pinger
type wrapConn0003_734a0c923476b7212c3571646d8223d0 interface {
	driver.ConnBeginTx
	driver.Execer
	driver.Pinger
}

// driver.ConnBeginTx|driver.Queryer|driver.QueryerContext
type wrapConn0003_580a0ede16c719a4c714bf892fdb4061 interface {
	driver.ConnBeginTx
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Queryer
type wrapConn0002_15a5fa7e0b58b2775996354e26f737f0 interface {
	driver.ConnBeginTx
	driver.Queryer
}

// driver.ConnPrepareContext|driver.Queryer
type wrapConn0002_2f1b86788f06ab353f2406475549da3a interface {
	driver.ConnPrepareContext
	driver.Queryer
}

// driver.Execer|driver.Queryer
type wrapConn0002_549f199f5c8eff6a17a1207eaca388f4 interface {
	driver.Execer
	driver.Queryer
}

// driver.ExecerContext|driver.Queryer
type wrapConn0002_863bdd334cd7bcaefd9d5085e2626cd9 interface {
	driver.ExecerContext
	driver.Queryer
}

// driver.NamedValueChecker|driver.Queryer
type wrapConn0002_d889ee74c9f7be4ca1b3d5f12965aecb interface {
	driver.NamedValueChecker
	driver.Queryer
}

// driver.Pinger|driver.Queryer
type wrapConn0002_36d783b278273e048cbac2a9db4fb974 interface {
	driver.Pinger
	driver.Queryer
}

// driver.ConnBeginTx|driver.ConnPrepareContext
type wrapConn0002_b7ff310fb7e22a7ccbb847bf26429dd9 interface {
	driver.ConnBeginTx
	driver.ConnPrepareContext
}

// driver.NamedValueChecker|driver.Pinger
type wrapConn0002_0921fb5e4c916e10370d5e0398435d4e interface {
	driver.NamedValueChecker
	driver.Pinger
}

// driver.ConnBeginTx|driver.QueryerContext
type wrapConn0002_21490f92aea80b01cad7a6a9b496de9d interface {
	driver.ConnBeginTx
	driver.QueryerContext
}

// driver.ConnPrepareContext|driver.QueryerContext
type wrapConn0002_e5960e40a7bbe1b5a19c38f71e0e6861 interface {
	driver.ConnPrepareContext
	driver.QueryerContext
}

// driver.Execer|driver.QueryerContext
type wrapConn0002_1edfe563a97fbd21d92eeed3f2776fd1 interface {
	driver.Execer
	driver.QueryerContext
}

// driver.ExecerContext|driver.QueryerContext
type wrapConn0002_9998c1796991ed1018aa6832ef10c20d interface {
	driver.ExecerContext
	driver.QueryerContext
}

// driver.NamedValueChecker|driver.QueryerContext
type wrapConn0002_608b01490fc5eeb54196421e44a1d294 interface {
	driver.NamedValueChecker
	driver.QueryerContext
}

// driver.Pinger|driver.QueryerContext
type wrapConn0002_fb43968160809b260946ae6f55b7fe56 interface {
	driver.Pinger
	driver.QueryerContext
}

// driver.SessionResetter|driver.Validator
type wrapConn0002_20cb62afffd12ef5a3dd67710522c25c interface {
	driver.SessionResetter
	driver.Validator
}

// driver.ExecerContext|driver.Pinger
type wrapConn0002_b0936442b3763765deb78c0f8945bb0a interface {
	driver.ExecerContext
	driver.Pinger
}

// driver.Queryer|driver.QueryerContext
type wrapConn0002_0ab84d337a68f7f945b7dd51c3322a35 interface {
	driver.Queryer
	driver.QueryerContext
}

// driver.ConnBeginTx|driver.Execer
type wrapConn0002_71ee22b4aad207f24e83c993b46f34e8 interface {
	driver.ConnBeginTx
	driver.Execer
}

// driver.QueryerContext|driver.Validator
type wrapConn0002_57475ba369865550e7711feefe48c5d4 interface {
	driver.QueryerContext
	driver.Validator
}

// driver.Execer|driver.Pinger
type wrapConn0002_270a739597d105c623663cb03227f1f9 interface {
	driver.Execer
	driver.Pinger
}

// driver.Queryer|driver.Validator
type wrapConn0002_592daa1271d02469a5a34fa382f5c410 interface {
	driver.Queryer
	driver.Validator
}

// driver.ConnBeginTx|driver.ExecerContext
type wrapConn0002_01ce0f63e6a6b920df90cbea6b6ede03 interface {
	driver.ConnBeginTx
	driver.ExecerContext
}

// driver.ConnBeginTx|driver.Pinger
type wrapConn0002_defc836eab973764559fef7a6eacab29 interface {
	driver.ConnBeginTx
	driver.Pinger
}

// driver.ConnBeginTx|driver.SessionResetter
type wrapConn0002_f21346c5fd58f22896f02b5b5ed3475e interface {
	driver.ConnBeginTx
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.SessionResetter
type wrapConn0002_1bf71fb4be83d1fefe01cdd4d5a7a547 interface {
	driver.ConnPrepareContext
	driver.SessionResetter
}

// driver.Execer|driver.SessionResetter
type wrapConn0002_3ca02c0eda0e299af94cd79f7897d2b4 interface {
	driver.Execer
	driver.SessionResetter
}

// driver.ExecerContext|driver.SessionResetter
type wrapConn0002_26864a2d0c8c5a605e497ec245ca9aa1 interface {
	driver.ExecerContext
	driver.SessionResetter
}

// driver.NamedValueChecker|driver.SessionResetter
type wrapConn0002_c4cbf403989171fe81ee42f1529c60fe interface {
	driver.NamedValueChecker
	driver.SessionResetter
}

// driver.ConnPrepareContext|driver.Execer
type wrapConn0002_a1e981f752f5cdd201647cde399eb7ea interface {
	driver.ConnPrepareContext
	driver.Execer
}

// driver.ExecerContext|driver.NamedValueChecker
type wrapConn0002_70a606d71e13e67187acc5d5056da2e1 interface {
	driver.ExecerContext
	driver.NamedValueChecker
}

// driver.Pinger|driver.Validator
type wrapConn0002_dea76879684fe95f8ca4bec7f7cec2d7 interface {
	driver.Pinger
	driver.Validator
}

// driver.NamedValueChecker|driver.Validator
type wrapConn0002_9743a8d2379220bf3bef64c2a633bd0d interface {
	driver.NamedValueChecker
	driver.Validator
}

// driver.Queryer|driver.SessionResetter
type wrapConn0002_40db605fa7c67cf6d53f67950cfe2a5a interface {
	driver.Queryer
	driver.SessionResetter
}

// driver.Execer|driver.NamedValueChecker
type wrapConn0002_7848e2d0c0ff74bcf6995d5ad7cbbd81 interface {
	driver.Execer
	driver.NamedValueChecker
}

// driver.ExecerContext|driver.Validator
type wrapConn0002_b1c510487a94a905aff86496f31078e6 interface {
	driver.ExecerContext
	driver.Validator
}

// driver.ConnPrepareContext|driver.NamedValueChecker
type wrapConn0002_35700b3f7c770fcd0bf83836eee99fec interface {
	driver.ConnPrepareContext
	driver.NamedValueChecker
}

// driver.ConnBeginTx|driver.NamedValueChecker
type wrapConn0002_1c65fb5735328219b65d13d4e717b79f interface {
	driver.ConnBeginTx
	driver.NamedValueChecker
}

// driver.Execer|driver.Validator
type wrapConn0002_45c298b388e6621ec2eb7a1ce46122bc interface {
	driver.Execer
	driver.Validator
}

// driver.ConnPrepareContext|driver.Validator
type wrapConn0002_739d949041a357f14dd543ae1c82d4b4 interface {
	driver.ConnPrepareContext
	driver.Validator
}

// driver.ConnBeginTx|driver.Validator
type wrapConn0002_7eb98af2c595d337fef40bcf00e367f5 interface {
	driver.ConnBeginTx
	driver.Validator
}

// driver.QueryerContext|driver.SessionResetter
type wrapConn0002_734bbe2891841d705d0385cb4817474b interface {
	driver.QueryerContext
	driver.SessionResetter
}

// driver.Execer|driver.ExecerContext
type wrapConn0002_a31656062b158c63465098f4075c1894 interface {
	driver.Execer
	driver.ExecerContext
}

// driver.ConnPrepareContext|driver.ExecerContext
type wrapConn0002_f2e0c829a13f0fc5e34c7487f6a5743b interface {
	driver.ConnPrepareContext
	driver.ExecerContext
}

// driver.ConnPrepareContext|driver.Pinger
type wrapConn0002_039f074dd8b7cd2fd02bcb0f3a89e29a interface {
	driver.ConnPrepareContext
	driver.Pinger
}

// driver.Pinger|driver.SessionResetter
type wrapConn0002_004b3f5abaa12288a30861da40988df6 interface {
	driver.Pinger
	driver.SessionResetter
}

// driver.ConnPrepareContext
type wrapConn0001_c3128ba459e963c7f07c0de8d58b775f interface {
	driver.ConnPrepareContext
}

// driver.QueryerContext
type wrapConn0001_7a610c95e7c563485e18c228cc42e7fb interface {
	driver.QueryerContext
}

// driver.Pinger
type wrapConn0001_f95e29e11298b900f8a05d72dcfe041b interface {
	driver.Pinger
}

// driver.SessionResetter
type wrapConn0001_70145ca1951f1697e67aa095f55b8305 interface {
	driver.SessionResetter
}

// driver.NamedValueChecker
type wrapConn0001_dd0fe84666fad76cbf0dee440e953d06 interface {
	driver.NamedValueChecker
}

// driver.Validator
type wrapConn0001_abc062388a1f7a907b94470050b8bc86 interface {
	driver.Validator
}

// driver.ConnBeginTx
type wrapConn0001_d815927e8e5ae23caed6edbf82ddb9a8 interface {
	driver.ConnBeginTx
}

// driver.ExecerContext
type wrapConn0001_f98304698b6793d84b0e7be539d135ce interface {
	driver.ExecerContext
}

// driver.Execer
type wrapConn0001_3380a887776f8c4760ee0d7b1fcdcde7 interface {
	driver.Execer
}

// driver.Queryer
type wrapConn0001_5f0836c749bf264b1429d45373cd34bc interface {
	driver.Queryer
}

func wrapStmt(stmt driver.Stmt, query string, opts Options) driver.Stmt {
	c := &wrapperStmt{stmt: stmt, query: query, opts: opts}
	if _, ok := stmt.(wrapStmt0004_c1de8b592171dd6c1ac443c171c04be4); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
			driver.ColumnConverter
			driver.NamedValueChecker
		}{c, c, c, c, c}
	}

	if _, ok := stmt.(wrapStmt0003_b9ce98633dc79f4a9c980cb302695741); ok {
		return struct {
			driver.Stmt
			driver.StmtQueryContext
			driver.ColumnConverter
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := stmt.(wrapStmt0003_99e9fa0a58cdfd18577f397c69971bfe); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.ColumnConverter
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := stmt.(wrapStmt0003_0856d47f81a59090013bce4a1ee5bdd5); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
			driver.NamedValueChecker
		}{c, c, c, c}
	}

	if _, ok := stmt.(wrapStmt0003_e6b11bef8ca1ccb0d5605a2c7d2a031c); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
			driver.ColumnConverter
		}{c, c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_f6c48d21623fccea572a7155d73bcb7e); ok {
		return struct {
			driver.Stmt
			driver.ColumnConverter
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_1ac318a0d0640d654e838a3d210e928f); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.StmtQueryContext
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_21dd080dbc547db7e0f1b0faf56aa2c3); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_85c9c41d2b9b8a7a8fcdecf72e140f4f); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
			driver.ColumnConverter
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_c4a6774a9c81bc05113b6cee4b0e4616); ok {
		return struct {
			driver.Stmt
			driver.StmtQueryContext
			driver.ColumnConverter
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0002_90dca7297ea2f892831234b29884308e); ok {
		return struct {
			driver.Stmt
			driver.StmtQueryContext
			driver.NamedValueChecker
		}{c, c, c}
	}

	if _, ok := stmt.(wrapStmt0001_6c51dc0f420a9aec8de551dde6d98c41); ok {
		return struct {
			driver.Stmt
			driver.ColumnConverter
		}{c, c}
	}

	if _, ok := stmt.(wrapStmt0001_be979ee6572ef74e683e3153d6b71c90); ok {
		return struct {
			driver.Stmt
			driver.StmtExecContext
		}{c, c}
	}

	if _, ok := stmt.(wrapStmt0001_3f3eb31555c2ec79768169cb473a84fe); ok {
		return struct {
			driver.Stmt
			driver.StmtQueryContext
		}{c, c}
	}

	if _, ok := stmt.(wrapStmt0001_dd0fe84666fad76cbf0dee440e953d06); ok {
		return struct {
			driver.Stmt
			driver.NamedValueChecker
		}{c, c}
	}

	return c
}

// driver.StmtExecContext|driver.StmtQueryContext|driver.ColumnConverter|driver.NamedValueChecker
type wrapStmt0004_c1de8b592171dd6c1ac443c171c04be4 interface {
	driver.StmtExecContext
	driver.StmtQueryContext
	driver.ColumnConverter
	driver.NamedValueChecker
}

// driver.StmtQueryContext|driver.ColumnConverter|driver.NamedValueChecker
type wrapStmt0003_b9ce98633dc79f4a9c980cb302695741 interface {
	driver.StmtQueryContext
	driver.ColumnConverter
	driver.NamedValueChecker
}

// driver.StmtExecContext|driver.ColumnConverter|driver.NamedValueChecker
type wrapStmt0003_99e9fa0a58cdfd18577f397c69971bfe interface {
	driver.StmtExecContext
	driver.ColumnConverter
	driver.NamedValueChecker
}

// driver.StmtExecContext|driver.StmtQueryContext|driver.NamedValueChecker
type wrapStmt0003_0856d47f81a59090013bce4a1ee5bdd5 interface {
	driver.StmtExecContext
	driver.StmtQueryContext
	driver.NamedValueChecker
}

// driver.StmtExecContext|driver.StmtQueryContext|driver.ColumnConverter
type wrapStmt0003_e6b11bef8ca1ccb0d5605a2c7d2a031c interface {
	driver.StmtExecContext
	driver.StmtQueryContext
	driver.ColumnConverter
}

// driver.ColumnConverter|driver.NamedValueChecker
type wrapStmt0002_f6c48d21623fccea572a7155d73bcb7e interface {
	driver.ColumnConverter
	driver.NamedValueChecker
}

// driver.StmtExecContext|driver.StmtQueryContext
type wrapStmt0002_1ac318a0d0640d654e838a3d210e928f interface {
	driver.StmtExecContext
	driver.StmtQueryContext
}

// driver.StmtExecContext|driver.NamedValueChecker
type wrapStmt0002_21dd080dbc547db7e0f1b0faf56aa2c3 interface {
	driver.StmtExecContext
	driver.NamedValueChecker
}

// driver.StmtExecContext|driver.ColumnConverter
type wrapStmt0002_85c9c41d2b9b8a7a8fcdecf72e140f4f interface {
	driver.StmtExecContext
	driver.ColumnConverter
}

// driver.StmtQueryContext|driver.ColumnConverter
type wrapStmt0002_c4a6774a9c81bc05113b6cee4b0e4616 interface {
	driver.StmtQueryContext
	driver.ColumnConverter
}

// driver.StmtQueryContext|driver.NamedValueChecker
type wrapStmt0002_90dca7297ea2f892831234b29884308e interface {
	driver.StmtQueryContext
	driver.NamedValueChecker
}

// driver.ColumnConverter
type wrapStmt0001_6c51dc0f420a9aec8de551dde6d98c41 interface {
	driver.ColumnConverter
}

// driver.StmtExecContext
type wrapStmt0001_be979ee6572ef74e683e3153d6b71c90 interface {
	driver.StmtExecContext
}

// driver.StmtQueryContext
type wrapStmt0001_3f3eb31555c2ec79768169cb473a84fe interface {
	driver.StmtQueryContext
}

// driver.NamedValueChecker
type wrapStmt0001_dd0fe84666fad76cbf0dee440e953d06 interface {
	driver.NamedValueChecker
}
